// What is this? It's an only-very-slightly modified copy of
// the starlark conversion package from https://github.com/cirruslabs/cirrus-cli/tree/8bdda497c6dcbf2d1c22bc13f8bdd046329ff874/pkg/larker
//
// The purpose of this package is to take a starlark input, parse it with the go starlark parser, resolve and load any additional
// starlark modules (since the .star files themselves can load others) and execute the 'main' function therein.
//
// Starlark itself has no notion of a 'module', this is a construct, based upon the made up by the cirrus-ci approach, but
// seems reasonable here.
//
// The intention of the main function is to act as a pure function which returns map or list struct which is then transformed into yaml
// by the 'larker' package.
package larker
