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

### Algorithm

Resolve is called from the [https://github.com/zond/godip/blob/master/judge/judge.go](judge).

Adjudicate is defined in each order type.

```
A depends on B: A = B
B depends on C: B = C
C = t

Resolve(A)
 Adjudicate(A)
  Resolve(B)
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
     -> t
    -> t
   -> t
  -> t
 -> t
-> t


A depends on B: A = B
B depends on C: B = C
C depends on A: C = A

Resolve(A)
 Adjudicate(A)
  Resolve(B)
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
      Resolve(A)
       # guess
      -> f
     -> f
    -> f
   -> f
  -> f
 -> f
 # guess of A was made
 Adjudicate(A)
  Resolve(B)
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
      Resolve(A)
       # guess
      -> t
     -> t
    -> t
   -> t
  -> t
 -> t
 # both guesses consistent
 BackupRule
-> X


A depends on B: A = B
B depends on C: B = C
C depends on A and D: C = A & D
D depends on A: D = !A

Resolve(A)
 Adjudicate(A)
  Resolve(B)
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
      Resolve(A)
       # guess
      -> f
      Resolve(D)
       Adjudicate(D)
        Resolve(A)
	 # already guessed
	-> f
       -> t
      -> t
     -> f
    -> f
   -> f
  -> f
 -> f
 # guess of A was made
 Adjudicate(A)
  Resolve(B)
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
      Resolve(A)
       # guess
      -> t
      Resolve(D)
       Adjudicate(D)
        Resolve(A)
	 # already guessed
	-> t
       -> f
      -> f
     -> f
    -> f
   -> f
  -> f
 -> f
 # only one guess consistent
-> f


A depends on B: A = B
B depends on C: B = C
C depends on A: C = !A

Resolve(A)
 Adjudicate(A)
  Resolve(B)
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
      Resolve(A)
       # guess
      -> f
     -> t
    -> t
   -> t
  -> t
 -> t
 # guess of A was made
 Adjudicate(A)
  Resolve(B)
   Adjudicate(B)
    Resolve(B)
     Adjudicate(C)
      Resolve(A)
       # guess
      -> t
     -> f
    -> f
   -> f
  -> f
 -> f
 # both guesses inconsistent
 BackupRule
-> X


A depends on B: A = B
B depends on C: B = C
C depends on B and A: C = B & A

Resolve(A)
 Adjudicate(A)
  Resolve(B)
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
      Resolve(B)
       # guess
      -> f
      Resolve(A)
       # guess
      -> f
     -> f
    -> f
   -> f
   # guess of B was made
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
      Resolve(B)
       # guess
      -> t
      Resolve(A)
       # already guessed
      -> f
     -> f
    -> f
   -> f
   # only one guess was consistent
  -> f
 -> f
 # guess of A was made
 Adjudicate(A)
  Resolve(B)
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
      Resolve(B)
       # guess
      -> f
      Resolve(A)
       # guess
      -> t
     -> f
    -> f
   -> f
   # guess of B was made
   Adjudicate(B)
    Resolve(C)
     Adjudicate(C)
      Resolve(B)
       # guess
      -> t
      Resolve(A)
       # already guessed
      -> t
     -> t
    -> t
   -> t
   # both guesses consistent
   BackupRule
  -> X
 -> X
 # if X is t, both guesses consistent: Y comes from BackupRule. if X is f, only one consistent, and Y is f
-> Y
```

