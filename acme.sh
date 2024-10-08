/Users/mac/wwwroot/acme.sh/acme.sh --renew -d lujiang.ren -d "*.lujiang.ren" --dns dns_ali --force
/Users/mac/wwwroot/acme.sh/acme.sh --renew -d 0001000.xyz -d "*.0001000.xyz" --dns dns_ali --force
/Users/mac/wwwroot/acme.sh/acme.sh --renew -d xingzhe.link -d "*.xingzhe.link" --dns dns_ali --force
/Users/mac/wwwroot/acme.sh/acme.sh --renew -d 88188.live -d "*.88188.live" --dns dns_ali --force
/Users/mac/wwwroot/acme.sh/acme.sh --renew -d 131.im -d "*.131.im" --dns dns_ali --force
/Users/mac/wwwroot/acme.sh/acme.sh --renew -d 818.gold -d "*.818.gold" --dns dns_ali --force
/Users/mac/wwwroot/acme.sh/acme.sh --renew -d 616.icu -d "*.616.icu" --dns dns_cf --force

cd /Users/mac/wwwroot/go-service/tcloud && ./tcloud cdn -d  adm.131.im -p /Users/mac/.acme.sh
cd /Users/mac/wwwroot/go-service/tcloud && ./tcloud cdn -d  r.131.im -p /Users/mac/.acme.sh
cd /Users/mac/wwwroot/go-service/tcloud && ./tcloud cdn -d  rs.131.im -p /Users/mac/.acme.sh
cd /Users/mac/wwwroot/go-service/tcloud && ./tcloud cdn -d  d.131.im -p /Users/mac/.acme.sh
cd /Users/mac/wwwroot/go-service/tcloud && ./tcloud cdn -d  static.131.im -p /Users/mac/.acme.sh
cd /Users/mac/wwwroot/go-service/tcloud && ./tcloud cdn -d  app.818.gold -p /Users/mac/.acme.sh
cd /Users/mac/wwwroot/go-service/tcloud && ./tcloud cdn -d  static.818.gold -p /Users/mac/.acme.sh
#~/go/src/tcloud/tcloud www.88188.live
#~/go/src/tcloud/tcloud www.0001000.xyz
#~/go/src/tcloud/tcloud www.xingzhe.link
#~/go/src/tcloud/tcloud www.lujiang.ren
#~/go/src/tcloud/tcloud www.818.gold
#~/go/src/tcloud/tcloud ai.818.gold
#~/go/src/tcloud/tcloud pwd.131.im
#~/go/src/tcloud/tcloud hpf.131.im
#~/go/src/tcloud/tcloud sec.131.im
#~/go/src/tcloud/tcloud mix.131.im
#~/go/src/tcloud/tcloud rs.131.im
#~/go/src/tcloud/tcloud oss.131.im
#~/go/src/tcloud/tcloud fang.131.im
#~/go/src/tcloud/tcloud mweb.131.im
#~/go/src/tcloud/tcloud glass.131.im