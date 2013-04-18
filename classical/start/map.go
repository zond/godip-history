package start

import (
  "github.com/zond/godip/classical/common"
  "github.com/zond/godip/graph"
)

func Graph() *graph.Graph {
  return graph.New().
    // nat
    Prov("nat").Conn("nrg").Conn("cly").Conn("lvp").Conn("iri").Conn("mid").Flag(common.Sea).
    // nrg
    Prov("nrg").Conn("nat").Conn("bar").Conn("nwy").Conn("nth").Conn("edi").Conn("cly").Flag(common.Sea).
    // bar
    Prov("bar").Conn("nrg").Conn("stp/nc").Conn("nwy").Flag(common.Sea).
    // stp/nc
    Prov("stp/nc").Conn("nrg").Conn("bar").Conn("nwy").Flag(common.Sea).
    // stp
    Prov("stp").Conn("fin").Conn("nwy").Conn("mos").Conn("lvn").Flag(common.Land).SC(common.Russia).
    // mos
    Prov("mos").Conn("stp").Conn("sev").Conn("ukr").Conn("war").Conn("lvn").Flag(common.Land).SC(common.Russia).
    // sev
    Prov("sev").Conn("ukr").Conn("mos").Conn("arm").Conn("bla").Conn("rum").Flag(common.Coast).SC(common.Russia).
    // arm
    Prov("arm").Conn("ank").Conn("bla").Conn("sev").Conn("syr").Conn("smy").Flag(common.Coast).
    // syr
    Prov("syr").Conn("eas").Conn("smy").Conn("arm").Flag(common.Coast).
    // eas
    Prov("eas").Conn("ion").Conn("aeg").Conn("smy").Conn("syr").Flag(common.Sea).
    // ion
    Prov("ion").Conn("tun").Conn("tys").Conn("nap").Conn("alb").Conn("gre").Conn("aeg").Conn("eas").Flag(common.Sea).
    // tun
    Prov("tun").Conn("naf").Conn("wes").Conn("tys").Conn("ion").Flag(common.Coast).SC(common.Neutral).
    // naf
    Prov("naf").Conn("mid").Conn("wes").Conn("tun").Flag(common.Coast).
    // mid
    Prov("mid").Conn("nat").Conn("iri").Conn("eng").Conn("bre").Conn("gas").Conn("spa/nc").Conn("por").Conn("spa/sc").Conn("naf").Flag(common.Sea).
    // iri
    Prov("iri").Conn("nat").Conn("lvp").Conn("wal").Conn("eng").Conn("mid").Flag(common.Sea).
    // lvp
    Prov("lvp").Conn("iri").Conn("nat").Conn("cly").Conn("edi").Conn("yor").Conn("wal").Flag(common.Coast).SC(common.England).
    // cly
    Prov("cly").Conn("nat").Conn("nrg").Conn("edi").Conn("lvp").Flag(common.Coast).
    // edi
    Prov("edi").Conn("cly").Conn("nrg").Conn("nth").Conn("yor").Conn("lvp").Flag(common.Coast).SC(common.England).
    // nth
    Prov("nth").Conn("edi").Conn("nrg").Conn("nwy").Conn("ska").Conn("den").Conn("hel").Conn("hel").Conn("hol").Conn("bel").Conn("lon").Conn("yor").Flag(common.Sea).
    // nwy
    Prov("nwy").Conn("nth").Conn("nrg").Conn("bar").Conn("stp/nc").Conn("stp").Conn("fin").Conn("swe").Conn("ska").Flag(common.Coast).SC(common.Neutral).
    // stp/sc
    Prov("stp/sc").Conn("bot").Conn("fin").Conn("lvn").Flag(common.Sea).
    // lvn
    Prov("lvn").Conn("bal").Conn("bot").Conn("stp/sc").Conn("mos").Conn("war").Conn("pru").Flag(common.Coast).
    // war
    Prov("war").Conn("sil").Conn("pru").Conn("lvn").Conn("mos").Conn("ukr").Conn("gal").Flag(common.Land).SC(common.Russia).
    // ukr
    Prov("ukr").Conn("war").Conn("mos").Conn("sev").Conn("rum").Conn("gal").Flag(common.Land).
    // bla
    Prov("bla").Conn("bul/ec").Conn("rum").Conn("sev").Conn("arm").Conn("ank").Conn("con").Flag(common.Sea).
    // ank
    Prov("ank").Conn("con").Conn("bla").Conn("arm").Conn("smy").Flag(common.Coast).SC(common.Turkey).
    // smy
    Prov("smy").Conn("aeg").Conn("con").Conn("ank").Conn("arm").Conn("syr").Conn("eas").Flag(common.Coast).SC(common.Turkey).
    // aeg
    Prov("aeg").Conn("ion").Conn("gre").Conn("bul/sc").Conn("con").Conn("smy").Flag(common.Sea).
    // gre
    Prov("gre").Conn("ion").Conn("alb").Conn("ser").Conn("bul").Conn("bul/sc").Conn("aeg").Flag(common.Coast).SC(common.Neutral).
    // nap
    Prov("nap").Conn("tys").Conn("rom").Conn("apu").Conn("ion").Flag(common.Coast).SC(common.Italy).
    // tys
    Prov("tys").Conn("wes").Conn("gol").Conn("tus").Conn("rom").Conn("nap").Conn("ion").Conn("tun").Flag(common.Sea).
    // wes
    Prov("wes").Conn("mid").Conn("spa/sc").Conn("gol").Conn("tys").Conn("tun").Conn("naf").Flag(common.Sea).
    // spa/sc
    Prov("spa/sc").Conn("mid").Conn("por").Conn("mar").Conn("gol").Conn("wes").Flag(common.Sea).
    // spa
    Prov("spa").Conn("por").Conn("gas").Conn("mar").Flag(common.Land).SC(common.Neutral).
    // spa/nc
    Prov("spa/nc").Conn("por").Conn("mid").Conn("gas").Flag(common.Sea).
    // por
    Prov("por").Conn("mid").Conn("spa/nc").Conn("spa").Conn("spa/sc").Flag(common.Coast).SC(common.Neutral).
    // gas
    Prov("gas").Conn("mid").Conn("bre").Conn("par").Conn("bur").Conn("mar").Conn("spa").Conn("spa/nc").Flag(common.Coast).
    // bre
    Prov("bre").Conn("mid").Conn("eng").Conn("pic").Conn("par").Conn("gas").Flag(common.Coast).SC(common.France).
    // eng
    Prov("eng").Conn("mid").Conn("iri").Conn("wal").Conn("lon").Conn("nth").Conn("bel").Conn("pic").Conn("bre").Flag(common.Sea).
    // wal
    Prov("wal").Conn("iri").Conn("lvp").Conn("yor").Conn("lon").Conn("eng").Flag(common.Coast).
    // yor
    Prov("yor").Conn("lvp").Conn("edi").Conn("nth").Conn("lon").Conn("wal").Flag(common.Coast).
    // ska
    Prov("ska").Conn("nth").Conn("nwy").Conn("swe").Conn("den").Flag(common.Sea).
    // swe
    Prov("swe").Conn("ska").Conn("nwy").Conn("fin").Conn("bot").Conn("bal").Conn("den").Flag(common.Coast).SC(common.Neutral).
    // fin
    Prov("fin").Conn("bot").Conn("swe").Conn("stp").Conn("stp/sc").Flag(common.Coast).
    // bot
    Prov("bot").Conn("swe").Conn("fin").Conn("stp/sc").Conn("lvn").Conn("bal").Flag(common.Sea).
    // bal
    Prov("bal").Conn("den").Conn("swe").Conn("bot").Conn("lvn").Conn("pru").Conn("ber").Conn("kie").Flag(common.Sea).
    // pru
    Prov("pru").Conn("ber").Conn("bal").Conn("lvn").Conn("war").Conn("sil").Flag(common.Coast).
    // sil
    Prov("sil").Conn("mun").Conn("ber").Conn("pru").Conn("war").Conn("gal").Conn("boh").Flag(common.Land).
    // gal
    Prov("gal").Conn("boh").Conn("sil").Conn("war").Conn("ukr").Conn("rum").Conn("bud").Conn("vie").Flag(common.Land).
    // rum
    Prov("rum").Conn("bud").Conn("gal").Conn("ukr").Conn("sev").Conn("bul/ec").Conn("bul").Conn("ser").Flag(common.Coast).SC(common.Neutral).
    // bul/ec
    Prov("bul/ec").Conn("rum").Conn("bla").Conn("con").Flag(common.Sea).
    // bul
    Prov("bul").Conn("ser").Conn("rum").Conn("con").Conn("gre").Flag(common.Land).SC(common.Neutral).
    // con
    Prov("con").Conn("bul/sc").Conn("bul").Conn("bul/ec").Conn("bla").Conn("ank").Conn("smy").Conn("aeg").Flag(common.Coast).SC(common.Turkey).
    // bul/sc
    Prov("bul/sc").Conn("gre").Conn("con").Conn("aeg").Flag(common.Sea).
    // ser
    Prov("ser").Conn("tri").Conn("bud").Conn("rum").Conn("bul").Conn("gre").Conn("alb").Flag(common.Land).SC(common.Neutral).
    // alb
    Prov("alb").Conn("adr").Conn("tri").Conn("ser").Conn("gre").Conn("ion").Flag(common.Coast).
    // adr
    Prov("adr").Conn("ven").Conn("tri").Conn("alb").Conn("ion").Conn("apu").Flag(common.Sea).
    // apu
    Prov("apu").Conn("rom").Conn("ven").Conn("adr").Conn("ion").Conn("nap").Flag(common.Coast).
    // rom
    Prov("rom").Conn("tys").Conn("tus").Conn("ven").Conn("apu").Conn("nap").Flag(common.Coast).SC(common.Italy).
    // tus
    Prov("tus").Conn("gol").Conn("pie").Conn("ven").Conn("rom").Conn("tys").Flag(common.Coast).
    // gol
    Prov("gol").Conn("spa/sc").Conn("mar").Conn("pie").Conn("tus").Conn("tys").Conn("wes").Flag(common.Sea).
    // mar
    Prov("mar").Conn("spa").Conn("gas").Conn("bur").Conn("pie").Conn("gol").Conn("spa/sc").Flag(common.Coast).SC(common.France).
    // bur
    Prov("bur").Conn("par").Conn("pic").Conn("bel").Conn("ruh").Conn("mun").Conn("mar").Conn("gas").Flag(common.Land).
    // par
    Prov("par").Conn("bre").Conn("pic").Conn("bur").Conn("gas").Flag(common.Land).SC(common.France).
    // pic
    Prov("pic").Conn("bre").Conn("eng").Conn("bel").Conn("bur").Conn("par").Flag(common.Coast).
    // lon
    Prov("lon").Conn("wal").Conn("yor").Conn("nth").Conn("eng").Flag(common.Coast).SC(common.England).
    // bel
    Prov("bel").Conn("pic").Conn("eng").Conn("nth").Conn("hol").Conn("ruh").Conn("bur").Flag(common.Coast).SC(common.Neutral).
    // hol
    Prov("hol").Conn("nth").Conn("hel").Conn("kie").Conn("ruh").Conn("bel").Flag(common.Coast).SC(common.Neutral).
    // hel
    Prov("hel").Conn("nth").Conn("den").Conn("kie").Conn("hol").Flag(common.Sea).
    // den
    Prov("den").Conn("hel").Conn("nth").Conn("ska").Conn("swe").Conn("bal").Conn("kie").Flag(common.Coast).SC(common.Neutral).
    // ber
    Prov("ber").Conn("kie").Conn("bal").Conn("pru").Conn("sil").Conn("mun").Flag(common.Coast).SC(common.Germany).
    // mun
    Prov("mun").Conn("bur").Conn("ruh").Conn("kie").Conn("ber").Conn("sil").Conn("boh").Conn("tyr").Flag(common.Land).SC(common.Germany).
    // boh
    Prov("boh").Conn("mun").Conn("sil").Conn("gal").Conn("vie").Conn("tyr").Flag(common.Land).
    // vie
    Prov("vie").Conn("tyr").Conn("boh").Conn("gal").Conn("bud").Conn("tri").Flag(common.Land).SC(common.Austria).
    // bud
    Prov("bud").Conn("tri").Conn("vie").Conn("gal").Conn("rum").Conn("ser").Flag(common.Land).SC(common.Austria).
    // tri
    Prov("tri").Conn("adr").Conn("ven").Conn("tyr").Conn("vie").Conn("bud").Conn("ser").Conn("alb").Flag(common.Coast).SC(common.Austria).
    // ven
    Prov("ven").Conn("tus").Conn("pie").Conn("tyr").Conn("tri").Conn("adr").Conn("apu").Conn("rom").Flag(common.Coast).SC(common.Italy).
    // pie
    Prov("pie").Conn("mar").Conn("tyr").Conn("ven").Conn("tus").Conn("gol").Flag(common.Coast).
    // ruh
    Prov("ruh").Conn("bel").Conn("hol").Conn("kie").Conn("mun").Conn("bur").Flag(common.Land).
    // tyr
    Prov("tyr").Conn("mun").Conn("boh").Conn("vie").Conn("tri").Conn("ven").Conn("pie").Flag(common.Land).
    // kie
    Prov("kie").Conn("hol").Conn("hel").Conn("den").Conn("bal").Conn("ber").Conn("mun").Conn("ruh").Flag(common.Coast).SC(common.Germany).
    Done()
}
