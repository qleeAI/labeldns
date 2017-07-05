package main

import (
	"fmt"
	"net"
    "flag"
    "strings"
    "math/rand"
	"github.com/miekg/dns"
)


var resolve_domain string

type dns_handler struct {
}

func (d dns_handler)ServeDNS(w dns.ResponseWriter, msg *dns.Msg) {
    qname := msg.Question[0].Name
    if resolve_domain != "" && resolve_domain != qname {
        fmt.Println("qname not match, qname: ", qname)
        return
    }

    //fmt.Println("msd: ", msg.MsgHdr)
    //fmt.Println("ques: ", msg.Question)
    //fmt.Println("remote addr: ", w.RemoteAddr())
    ip := ""
    for i := 0; i < 4; i++ {
        num := rand.Intn(256)
        ip = ip + fmt.Sprintf("%d", num)
        if i != 3 {
            ip = ip + "."
        }
    }

    r, err  := dns.NewRR(qname + " A " + ip)
    if err != nil {
        fmt.Println("new rr failed")
        return
    }

    resp := msg.SetReply(msg)
    //resp.Question = make
    resp.Question = make([]dns.Question, 1)
    resp.Question[0] = dns.Question{qname, dns.TypeA, dns.ClassINET}
    resp.Answer = make([]dns.RR, 1)
    resp.Answer[0] = r
    err = w.WriteMsg(resp)
    if err != nil {
        fmt.Println("writemsg failed, err: ", err)
        return
    }
    fmt.Println("remote_addr: ", w.RemoteAddr(), ", ip: ", ip)
    //fmt.Println("resp: ", resp)
    return
}

func main() {

    port := flag.String("port", "53", "bind port")
    ip := flag.String("host", "", "bind ip address")
    domain := flag.String("domain", "", "resolve domain")
    flag.Parse()

    if *domain != "" && !strings.HasSuffix(*domain, ".") {
        resolve_domain = *domain + "."
    } else {
        resolve_domain = *domain
    }
    //fmt.Println("ip: ", *ip, ", port: ", *port, ", domain: ", *domain)

    ln, err := net.ListenPacket("udp", *ip + ":" + *port)
    if err != nil {
        fmt.Println(err)
        return
	}

    var d dns_handler
    dns.ActivateAndServe(nil, ln, d)

    fmt.Println("end")
}
