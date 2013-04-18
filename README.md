godip
=====

A dippy judge in Go.

### Policy

This is, currently, my thoughts on what must be true in all variants this judge shall support.

* There is a Map with Provinces.
 * Each Province can contain at most one Unit.
 * Each Province can contain at most one SupplyCenter.
* Each Province has one or more SubProvinces.
 * Each SubProvince has attributes, such as Sea, Land or Coast.
 * Each SubProvince has connections to other SubProvinces.
* There are Units with UnitType and Nationality.
 * Each unit is in one SubProvince.
* There are SupplyCenters with Nationality.
 * Each SupplyCenter is in one Province.
* There are Phases with Year, Season and PhaseType.
* Orders can vary greatly:
 * They can be valid only for certain Years, PhaseTypes or Seasons.
 * They can be valid only for certain UnitTypes.
 * They can be valid only for certain Map environments.
 * They can be valid only when certain other Orders are valid.
 * They can be valid only when certain Units are present.
