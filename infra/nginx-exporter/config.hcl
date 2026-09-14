listen {
  address = "0.0.0.0"
  port = 4040
}

namespace "korp_nginx" {
  source {
    syslog {
      listen_address = "udp://127.0.0.1:5531"
      format = "rfc3164"
      tags = ["korp"]
    }
  }

  format = "\"$request\" $status $body_bytes_sent $request_time"
  histogram_buckets = [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
}
