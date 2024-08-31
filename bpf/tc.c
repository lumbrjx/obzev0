#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/pkt_cls.h>
#include <linux/tcp.h>
#include <linux/udp.h>

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

struct {
  __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
  __uint(key_size, sizeof(int));
  __uint(value_size, sizeof(int));
  __uint(max_entries, 1024);
} events SEC(".maps");

struct event {
  __u32 src_ip;
  __u32 dst_ip;
  __u16 src_port;
  __u16 dst_port;
  __u8 protocol;
  __u8 direction;
  __u8 tcp_flags; 
};

// to avoid duplication :>
static __always_inline int process_packet(struct __sk_buff *skb,
                                          unsigned char direction) {
  void *data = (void *)(unsigned long)skb->data;
  void *data_end = (void *)(unsigned long)skb->data_end;

  // eht header
  struct ethhdr *eth = data;
  if ((void *)(eth + 1) > data_end)
    return TC_ACT_SHOT;
  if (eth->h_proto != bpf_htons(ETH_P_IP))
    return TC_ACT_OK;

  // ip header
  struct iphdr *ip = (struct iphdr *)(eth + 1);
  if ((void *)(ip + 1) > data_end)
    return TC_ACT_SHOT;

  // event creation
  struct event e = {0};
  e.src_ip = ip->saddr;
  e.dst_ip = ip->daddr;
  e.protocol = ip->protocol;
  e.direction = direction;

  if (ip->protocol == IPPROTO_TCP) {
    struct tcphdr *tcp = (struct tcphdr *)(ip + 1);
    if ((void *)(tcp + 1) > data_end)
      return TC_ACT_SHOT;

    e.src_port = bpf_ntohs(tcp->source);
    e.dst_port = bpf_ntohs(tcp->dest);
    e.tcp_flags = tcp->fin | (tcp->syn << 1) | (tcp->rst << 2) |
                  (tcp->psh << 3) | (tcp->ack << 4) | (tcp->urg << 5) |
                  (tcp->ece << 6) | (tcp->cwr << 7);
  } else if (ip->protocol == IPPROTO_UDP) {
    struct udphdr *udp = (struct udphdr *)(ip + 1);
    if ((void *)(udp + 1) > data_end)
      return TC_ACT_SHOT;

    e.src_port = bpf_ntohs(udp->source);
    e.dst_port = bpf_ntohs(udp->dest);
  } else {
    e.src_port = 0;
    e.dst_port = 0;
    e.tcp_flags = 0;
  }
  // outputing the data via a perf event array map
  bpf_perf_event_output(skb, &events, BPF_F_CURRENT_CPU, &e, sizeof(e));

  return TC_ACT_OK;
}

SEC("tc")
int tc_ingress(struct __sk_buff *skb) { return process_packet(skb, 0); }

SEC("tc")
int tc_egress(struct __sk_buff *skb) { return process_packet(skb, 1); }

char _license[] SEC("license") = "GPL";
