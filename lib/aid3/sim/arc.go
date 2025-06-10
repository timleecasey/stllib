package sim

// Format 1: (G17) G02/03 X__ Y__ I__ J__ F__
// (G18) G02/03 X__ Z__ I__ K__ F__
// (G19) G02/03 Y__ Z__ J__ K__ F__
//
// I,J,K specify the current plane analogous to X,Y,Z
// Only one of the three planes is used at any one time
// X/Y/Z I/J -> xy plane and z == k
// X/Y/Z I/K -> xz plane and y == j
// X/Y/Z J/K -> yz plane and x == i
//

// Format 2: (G17)G02/03 X__ Y__ R__ F__
// (G18)G02/03 X__ Z__ R__ F__
// (G19)G02/03 Y__ Z__ R__ F__
//
// Same game.
// X/Y/Z X/Y/R xy plane z is constant
// X/Y/Z X/Z/R xz plane y is constant
// X/Y/Z Y/Z/R yz plane x is constant
//
