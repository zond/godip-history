package classical

import (
  "github.com/zond/godip/graph"
)

const (
  Sea = 1 << iota
  Land
)

const (
  Coast = Sea | Land
  SC    = "SC"

  England = "England"
  France  = "France"
  Germany = "Germany"
  Russia  = "Russia"
  Austria = "Austria"
  Italy   = "Italy"
  Turkey  = "Turkey"
  Neutral = "Neutral"
)

func Graph() *graph.Graph {
  return graph.New().
    // nat
    Node("nat").Con("nrg").Con("cly").Con("lvp").Con("iri").Con("mid").Flag(Sea).
    // nrg
    Node("nrg").Con("nat").Con("bar").Con("nwy").Con("nth").Con("edi").Con("cly").Flag(Sea).
    // bar
    Node("bar").Con("nrg").Conn("stp", "nc").Con("nwy").Flag(Sea).
    // stp/nc
    Node("stp").Sub("nc").Con("nrg").Con("bar").Con("nwy").Flag(Sea).
    // stp
    Node("stp").Con("fin").Con("nwy").Con("mos").Con("lvn").Flag(Land).Attr(SC, Russia).
    // mos
    Node("mos").Con("stp").Con("sev").Con("ukr").Con("war").Con("lvn").Flag(Land).Attr(SC, Russia).
    // sev
    Node("sev").Con("ukr").Con("mos").Con("arm").Con("bla").Con("rum").Flag(Coast).Attr(SC, Russia).
    // arm
    Node("arm").Con("ank").Con("bla").Con("sev").Con("syr").Con("smy").Flag(Coast).
    // syr
    Node("syr").Con("eas").Con("smy").Con("arm").Flag(Coast).
    // eas
    Node("eas").Con("ion").Con("aeg").Con("smy").Con("syr").Flag(Sea).
    // ion
    Node("ion").Con("tun").Con("tys").Con("nap").Con("alb").Con("gre").Con("aeg").Con("eas").Flag(Sea).
    // tun
    Node("tun").Con("naf").Con("wes").Con("tys").Con("ion").Flag(Coast).Attr(SC, Neutral).
    // naf
    Node("naf").Con("mid").Con("wes").Con("tun").Flag(Coast).
    // mid
    Node("mid").Con("nat").Con("iri").Con("eng").Con("bre").Con("gas").Conn("spa", "nc").Con("por").Conn("spa", "sc").Con("naf").Flag(Sea).
    // iri
    Node("iri").Con("nat").Con("lvp").Con("wal").Con("eng").Con("mid").Flag(Sea).
    // lvp
    Node("lvp").Con("iri").Con("nat").Con("cly").Con("edi").Con("yor").Con("wal").Flag(Coast).Attr(SC, England).
    // cly
    Node("cly").Con("nat").Con("nrg").Con("edi").Con("lvp").Flag(Coast).
    // edi
    Node("edi").Con("cly").Con("nrg").Con("nth").Con("yor").Con("lvp").Flag(Coast).Attr(SC, England).
    // nth
    Node("nth").Con("edi").Con("nrg").Con("nwy").Con("ska").Con("den").Con("hel").Con("hel").Con("hol").Con("bel").Con("lon").Con("yor").Flag(Sea).
    // nwy
    Node("nwy").Con("nth").Con("nrg").Con("bar").Conn("stp", "nc").Con("stp").Con("fin").Con("swe").Con("ska").Flag(Coast).Attr(SC, Neutral).
    // stp/sc
    Node("stp").Sub("sc").Con("bot").Con("fin").Con("lvn").Flag(Sea).
    // lvn
    Node("lvn").Con("bal").Con("bot").Conn("stp", "sc").Con("mos").Con("war").Con("pru").Flag(Coast).
    // war
    Node("war").Con("sil").Con("pru").Con("lvn").Con("mos").Con("ukr").Con("gal").Flag(Land).Attr(SC, Russia).
    // ukr
    Node("ukr").Con("war").Con("mos").Con("sev").Con("rum").Con("gal").Flag(Land).
    // bla
    Node("bla").Conn("bul", "ec").Con("rum").Con("sev").Con("arm").Con("ank").Con("con").Flag(Sea).
    // ank
    Node("ank").Con("con").Con("bla").Con("arm").Con("smy").Flag(Coast).Attr(SC, Turkey).
    // smy
    Node("smy").Con("aeg").Con("con").Con("ank").Con("arm").Con("syr").Con("eas").Flag(Coast).Attr(SC, Turkey).
    // aeg
    Node("aeg").Con("ion").Con("gre").Conn("bul", "sc").Con("con").Con("smy").Flag(Sea).
    // gre
    Node("gre").Con("ion").Con("alb").Con("ser").Con("bul").Conn("bul", "sc").Con("aeg").Flag(Coast).Attr(SC, Neutral).
    // nap
    Node("nap").Con("tys").Con("rom").Con("apu").Con("ion").Flag(Coast).Attr(SC, Italy).
    // tys
    Node("tys").Con("wes").Con("gol").Con("tus").Con("rom").Con("nap").Con("ion").Con("tun").Flag(Sea).
    // wes
    Node("wes").Con("mid").Conn("spa", "sc").Con("gol").Con("tys").Con("tun").Con("naf").Flag(Sea).
    // spa/sc
    Node("spa").Sub("sc").Con("mid").Con("por").Con("mar").Con("gol").Con("wes").Flag(Sea).
    // spa
    Node("spa").Con("por").Con("gas").Con("mar").Flag(Land).Attr(SC, Neutral).
    // spa/nc
    Node("spa").Sub("nc").Con("por").Con("mid").Con("gas").Flag(Sea).
    // por
    Node("por").Con("mid").Conn("spa", "nc").Con("spa").Conn("spa", "sc").Flag(Coast).Attr(SC, Neutral).
    // gas
    Node("gas").Con("mid").Con("bre").Con("par").Con("bur").Con("mar").Con("spa").Conn("spa", "nc").Flag(Coast).
    // bre
    Node("bre").Con("mid").Con("eng").Con("pic").Con("par").Con("gas").Flag(Coast).Attr(SC, France).
    // eng
    Node("eng").Con("mid").Con("iri").Con("wal").Con("lon").Con("nth").Con("bel").Con("pic").Con("bre").Flag(Sea).
    // wal
    Node("wal").Con("iri").Con("lvp").Con("yor").Con("lon").Con("eng").Flag(Coast).
    // yor
    Node("yor").Con("lvp").Con("edi").Con("nth").Con("lon").Con("wal").Flag(Coast).
    // ska
    Node("ska").Con("nth").Con("nwy").Con("swe").Con("den").Flag(Sea).
    // swe
    Node("swe").Con("ska").Con("nwy").Con("fin").Con("bot").Con("bal").Con("den").Flag(Coast).Attr(SC, Neutral).
    // fin
    Node("fin").Con("bot").Con("swe").Con("stp").Conn("stp", "sc").Flag(Coast).
    // bot
    Node("bot").Con("swe").Con("fin").Conn("stp", "sc").Con("lvn").Con("bal").Flag(Sea).
    // bal
    Node("bal").Con("den").Con("swe").Con("bot").Con("lvn").Con("pru").Con("ber").Con("kie").Flag(Sea).
    // pru
    Node("pru").Con("ber").Con("bal").Con("lvn").Con("war").Con("sil").Flag(Coast).
    // sil
    Node("sil").Con("mun").Con("ber").Con("pru").Con("war").Con("gal").Con("boh").Flag(Land).
    // gal
    Node("gal").Con("boh").Con("sil").Con("war").Con("ukr").Con("rum").Con("bud").Con("vie").Flag(Land).
    // rum
    Node("rum").Con("bud").Con("gal").Con("ukr").Con("sev").Conn("bul", "ec").Con("bul").Con("ser").Flag(Coast).Attr(SC, Neutral).
    // bul/ec
    Node("bul").Sub("ec").Con("rum").Con("bla").Con("con").Flag(Sea).
    // bul
    Node("bul").Con("ser").Con("rum").Con("con").Con("gre").Flag(Land).Attr(SC, Neutral).
    // con
    Node("con").Conn("bul", "sc").Con("bul").Conn("bul", "ec").Con("bla").Con("ank").Con("smy").Con("aeg").Flag(Coast).Attr(SC, Turkey).
    // bul/sc
    Node("bul").Sub("sc").Con("gre").Con("con").Con("aeg").Flag(Sea).
    // ser
    Node("ser").Con("tri").Con("bud").Con("rum").Con("bul").Con("gre").Con("alb").Flag(Land).Attr(SC, Neutral).
    // alb
    Node("alb").Con("adr").Con("tri").Con("ser").Con("gre").Con("ion").Flag(Coast).
    // adr
    Node("adr").Con("ven").Con("tri").Con("alb").Con("ion").Con("apu").Flag(Sea).
    // apu
    Node("apu").Con("rom").Con("ven").Con("adr").Con("ion").Con("nap").Flag(Coast).
    // rom
    Node("rom").Con("tys").Con("tus").Con("ven").Con("apu").Con("nap").Flag(Coast).Attr(SC, Italy).
    // tus
    Node("tus").Con("gol").Con("pie").Con("ven").Con("rom").Con("tys").Flag(Coast).
    // gol
    Node("gol").Conn("spa", "sc").Con("mar").Con("pie").Con("tus").Con("tys").Con("wes").Flag(Sea).
    // mar
    Node("mar").Con("spa").Con("gas").Con("bur").Con("pie").Con("gol").Conn("spa", "sc").Flag(Coast).Attr(SC, France).
    // bur
    Node("bur").Con("par").Con("pic").Con("bel").Con("ruh").Con("mun").Con("mar").Con("gas").Flag(Land).
    // par
    Node("par").Con("bre").Con("pic").Con("bur").Con("gas").Flag(Land).Attr(SC, France).
    // pic
    Node("pic").Con("bre").Con("eng").Con("bel").Con("bur").Con("par").Flag(Coast).
    // lon
    Node("lon").Con("wal").Con("yor").Con("nth").Con("eng").Flag(Coast).Attr(SC, England).
    // bel
    Node("bel").Con("pic").Con("eng").Con("nth").Con("hol").Con("ruh").Con("bur").Flag(Coast).Attr(SC, Neutral).
    // hol
    Node("hol").Con("nth").Con("hel").Con("kie").Con("ruh").Con("bel").Flag(Coast).Attr(SC, Neutral).
    // hel
    Node("hel").Con("nth").Con("den").Con("kie").Con("hol").Flag(Sea).
    // den
    Node("den").Con("hel").Con("nth").Con("ska").Con("swe").Con("bal").Con("kie").Flag(Coast).Attr(SC, Neutral).
    // ber
    Node("ber").Con("kie").Con("bal").Con("pru").Con("sil").Con("mun").Flag(Coast).Attr(SC, Germany).
    // mun
    Node("mun").Con("bur").Con("ruh").Con("kie").Con("ber").Con("sil").Con("boh").Con("tyr").Flag(Land).Attr(SC, Germany).
    // boh
    Node("boh").Con("mun").Con("sil").Con("gal").Con("vie").Con("tyr").Flag(Land).
    // vie
    Node("vie").Con("tyr").Con("boh").Con("gal").Con("bud").Con("tri").Flag(Land).Attr(SC, Austria).
    // bud
    Node("bud").Con("tri").Con("vie").Con("gal").Con("rum").Con("ser").Flag(Land).Attr(SC, Austria).
    // tri
    Node("tri").Con("adr").Con("ven").Con("tyr").Con("vie").Con("bud").Con("ser").Con("alb").Flag(Coast).Attr(SC, Austria).
    // ven
    Node("ven").Con("tus").Con("pie").Con("tyr").Con("tri").Con("adr").Con("apu").Con("rom").Flag(Coast).Attr(SC, Italy).
    // pie
    Node("pie").Con("mar").Con("tyr").Con("ven").Con("tus").Con("gol").Flag(Coast).
    // ruh
    Node("ruh").Con("bel").Con("hol").Con("kie").Con("mun").Con("bur").Flag(Land).
    // tyr
    Node("tyr").Con("mun").Con("boh").Con("vie").Con("tri").Con("ven").Con("pie").Flag(Land).
    // kie
    Node("kie").Con("hol").Con("hel").Con("den").Con("bal").Con("ber").Con("mun").Con("ruh").Flag(Coast).Attr(SC, Germany).
    Done()
}
