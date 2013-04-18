godip
=====

A dippy judge in Go.

### Policy

This is, currently, my thoughts on what must be true in all variants this judge shall support.

* There is a map with Provinces.
 * Each Province has one or more SubProvinces.
 * Each Province can contain at most one Unit.
 * Each Province can contain at most one SupplyCenter.
 * Each SubProvince has connections to other SubProvinces.
 * Each SubProvince has attributes, such as Sea, Land, Coast.
* There are Units with UnitType and Nationality.
 * Each unit is in one SubProvince.
* There are SupplyCenters with Nationality.
 * Each SupplyCenter is in one Province.
* There are Phases with Year, Season and Type.
