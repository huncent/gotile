# gotile
basically a better in every conceivable way version of tile_reduce https://github.com/murphy214/tile_reduce

### INFO
This is a project I built to replace some of the core functionality I built into tile reduce, basically it takes the same basic concept, and implements it under a clean structure (geojson) instead of fractal structs for each potential geometry. Essentially now it can take raw geojson collections or ingest postgis database under the condition that there all single feature geometries (i.e. it currently doesn't support multiple geometries) this should be hard to fix I just figured I'd get the base 3 done first.

Beyond massively reducing the messiness of the code (streamlining all geometries into one process essentially) it also reconfigured key parts and has some weird experiments with recursion currently in attempts to reduce memory footprint. (which works pretty good) However the code currently has like 3 implementation of each pipeline just to test out each one. 

### TODO
Currently there doesn't exist a complete implementation of well known binary in golang without tons of shit to modify to get it to play nice with my structures. Therefore it currently parses the string geometry which is disgusting. I usually use a raw string field containing the geometry then use the json module to bring it into memory, but postgis implementations was kind of the goal. 

I can implement a wkb decoder 100% I just don't feel like getting into it currently. 

