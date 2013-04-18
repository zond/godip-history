godip
=====

A dippy judge in Go.

### Policy

* There are Units with UnitType and Nationality.
* There are SupplyCenters with Nationality.
* There are Phases with Year, Season and Type.
* There is a map, a Graph, with Provinces.
 * Each Province can have several sub provinces, ie coasts.
 * Each Province can hold one unit.
 * Each coast (including the "empty coast") has connections to other coasts.
* There are variants, and they differ in regards to
 * Available Orders.
 * Available UnitTypes.
 * Phase progression.
