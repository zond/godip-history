package graph

var Classical = New().
  // nat
  Node("nat").Con("nrg").Con("cly").Con("lvp").Con("iri").Con("mid").
  // nrg
  Node("nrg").Con("nat").Con("bar").Con("nwy").Con("nth").Con("edi").Con("cly").
  // bar
  Node("bar").Con("nrg").Conn("stp", "nc").Con("nwy").
  // stp/nc
  Node("stp").Sub("nc").Con("nrg").Con("bar").Con("nwy").
  // stp
  Node("stp").Con("fin").Con("nwy").Con("mos").Con("lvn").
  // mos
  Node("mos").Con("stp").Con("sev").Con("ukr").Con("war").Con("lvn").
  // sev
  Node("sev").Con("ukr").Con("mos").Con("arm").Con("bla").Con("rum").
  // arm
  Node("arm").Con("ank").Con("bla").Con("sev").Con("syr").Con("smy").
  // syr
  Node("syr").Con("eas").Con("smy").Con("arm").
  // eas
  Node("eas").Con("ion").Con("aeg").Con("smy").Con("syr").
  // ion
  Node("ion").Con("tun").Con("tys").Con("nap").Con("alb").Con("gre").Con("aeg").Con("eas").
  // tun
  Node("tun").Con("naf").Con("wes").Con("tys").Con("ion").
  // naf
  Node("naf").Con("mid").Con("wes").Con("tun").
  // mid
  Node("mid").Con("nat").Con("iri").Con("eng").Con("bre").Con("gas").Conn("spa", "nc").Con("por").Conn("spa", "sc").Con("naf").
  // iri
  Node("iri").Con("nat").Con("lvp").Con("wal").Con("eng").Con("mid").
  // lvp
  Node("lvp").Con("iri").Con("nat").Con("cly").Con("edi").Con("yor").Con("wal").
  // cly
  Node("cly").Con("nat").Con("nrg").Con("edi").Con("lvp").
  // edi
  Node("edi").Con("cly").Con("nrg").Con("nth").Con("yor").Con("lvp").
  // nth
  Node("nth").Con("edi").Con("nrg").Con("nwy").Con("ska").Con("den").Con("hel").Con("hel").Con("hol").Con("bel").Con("lon").Con("yor").
  // nwy
  Node("nwy").Con("nth").Con("nrg").Con("bar").Conn("stp", "nc").Con("stp").Con("fin").Con("swe").Con("ska").
  // stp/sc
  Node("stp").Sub("sc").Con("bot").Con("fin").Con("lvn").
  // lvn
  Node("lvn").Con("bal").Con("bot").Conn("stp", "sc").Con("mos").Con("war").Con("pru").
  // war
  Node("war").Con("sil").Con("pru").Con("lvn").Con("mos").Con("ukr").Con("gal").
  // ukr
  Node("ukr").Con("war").Con("mos").Con("sev").Con("rum").Con("gal").
  // bla
  Node("bla").Conn("bul", "ec").Con("rum").Con("sev").Con("arm").Con("ank").Con("con").
  // ank
  Node("ank").Con("con").Con("bla").Con("arm").Con("smy").
  // smy
  Node("smy").Con("aeg").Con("con").Con("ank").Con("arm").Con("syr").Con("eas").
  // aeg
  Node("aeg").Con("ion").Con("gre").Conn("bul", "sc").Con("con").Con("smy").
  // gre
  Node("gre").Con("ion").Con("alb").Con("ser").Con("bul").Conn("bul", "sc").Con("aeg").
  // nap
  Node("nap").Con("tys").Con("rom").Con("apu").Con("ion").
  // tys
  Node("tys").Con("wes").Con("gol").Con("tus").Con("rom").Con("nap").Con("ion").Con("tun").
  // wes
  Node("wes").Con("mid").Conn("spa", "sc").Con("gol").Con("tys").Con("tun").Con("naf").
  // spa/sc
  Node("spa").Sub("sc").Con("mid").Con("por").Con("mar").Con("gol").Con("wes").
  // spa
  Node("spa").Con("por").Con("gas").Con("mar").
  // por
  Node("por").Con("mid").Conn("spa", "nc").Con("spa").Conn("spa", "sc").
  // gas
  Node("gas").Con("mid").Con("bre").Con("par").Con("bur").Con("mar").Con("spa").Conn("spa", "nc").
  // bre
  Node("bre").Con("mid").Con("eng").Con("pic").Con("par").Con("gas").
  // eng
  Node("eng").Con("mid").Con("iri").Con("wal").Con("lon").Con("nth").Con("bel").Con("pic").Con("bre").
  // wal
  Node("wal").Con("iri").Con("lvp").Con("yor").Con("lon").Con("eng").
  // yor
  Node("yor").Con("lvp").Con("edi").Con("nth").Con("lon").Con("wal").
  // ska
  Node("ska").Con("nth").Con("nwy").Con("swe").Con("den").
  // swe
  Node("swe").Con("ska").Con("nwy").Con("fin").Con("bot").Con("bal").Con("den").
  // fin
  Node("fin").Con("bot").Con("swe").Con("stp").Conn("stp", "sc").
  // bot
  Node("bot").Con("swe").Con("fin").Conn("stp", "sc").Con("lvn").Con("bal").
  // bal
  Node("bal").Con("den").Con("swe").Con("bot").Con("lvn").Con("pru").Con("ber").Con("kie").
  // pru
  Node("pru").Con("ber").Con("bal").Con("lvn").Con("war").Con("sil").
  // sil
  Node("sil").Con("mun").Con("ber").Con("pru").Con("war").Con("gal").Con("boh").
  // gal
  Node("gal").Con("boh").Con("sil").Con("war").Con("ukr").Con("rum").Con("bud").Con("vie").
  // rum
  Node("rum").Con("bud").Con("gal").Con("ukr").Con("sev").Conn("bul", "ec").Con("bul").Con("ser").
  // bul/ec
  Node("bul").Sub("ec").Con("rum").Con("bla").Con("con").
  // bul
  Node("bul").Con("ser").Con("rum").Con("con").Con("gre").
  // con
  Node("con").Conn("bul", "sc").Con("bul").Conn("bul", "ec").Con("bla").Con("ank").Con("smy").Con("aeg").
  // bul/sc
  Node("bul").Sub("sc").Con("gre").Con("con").Con("aeg").
  // ser
  Node("ser").Con("tri").Con("bud").Con("rum").Con("bul").Con("gre").Con("alb").
  // alb
  Node("alb").Con("adr").Con("tri").Con("ser").Con("gre").Con("ion").
  // adr
  Node("adr").Con("ven").Con("tri").Con("alb").Con("ion").Con("apu").
  // apu
  Node("apu").Con("rom").Con("ven").Con("adr").Con("ion").Con("nap").
  // rom
  Node("rom").Con("tys").Con("tus").Con("ven").Con("apu").Con("nap").
  // tus
  Node("tus").Con("gol").Con("pie").Con("ven").Con("rom").Con("tys").
  // gol
  Node("gol").Conn("spa", "sc").Con("mar").Con("pie").Con("tus").Con("tys").Con("wes").
  // mar
  Node("mar").Con("spa").Con("gas").Con("bur").Con("pie").Con("gol").Conn("spa", "sc").
  // bur
  Node("bur").Con("par").Con("pic").Con("bel").Con("ruh").Con("mun").Con("mar").Con("gas").
  // par
  Node("par").Con("bre").Con("pic").Con("bur").Con("gas").
  // pic
  Node("pic").Con("bre").Con("eng").Con("bel").Con("bur").Con("par").
  // lon
  Node("lon").Con("wal").Con("yor").Con("nth").Con("eng").
  // bel
  Node("bel").Con("pic").Con("eng").Con("nth").Con("hol").Con("ruh").Con("bur").
  // hol
  Node("hol").Con("nth").Con("hel").Con("kie").Con("ruh").Con("bel").
  Done()
