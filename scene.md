# Scene builder


Scene file format:

```yaml
fixtures: # can have any number of lights
- type: pointlight 
  position: [ x, y, z ] # floats
  color: [ r, g, b ] # floats
materials: # A material dictionary that can be used in objects below
- name:   # name of the material
  preset: # Any item in the material cache that appears above this item, or "glass" or "default"
  pattern: solid | gradient | ring | checker
  colors: # Array of [r, g, b] colors to be used in the pattern. 1 color for "solid", 2 colors for the rest
  transform: # optional section, defaults to identity
  - type : identity | translate | scale | rotatex | rotatey | rotatez | shear
    params: # an array of floats that matches the transform type
            # identity - no params
            # translate - [ x, y, z ] floats
            # scale - [ x, y, z ] floats
            # rotate x,y,z - [ radians ] 
            # shear [ xy, xz, yx, yz, zx, zy ] floats 
  # all of the following material parameters are optional. The default is either the one written, or the one provided by the preset (if used)
  ambient: # float in the inclusive range [0,1]. Defaults to 0.1
  diffuse:  # float in the inclusive range [0,1]. Defaults to 0.9
  specular: # float in the inclusive range [0,1]. Defaults to 0/9
  shininess: # float in the inclusive range [0, inf ]. Defaults to 200
  reflective: # float in the inclusive range [0,1]. Defaults to 0.0
  transparency:  # float in the inclusive range [0,1]: Defaults to 0.0
  refractiveIndex: # float in the inclusive range [0,inf ]: Defaults to 1.0 
camera:
  hsize: # horizontal size of the rendered image
  vsize: # vertical size of the rendered image
  fieldOfView: # radians of the field of view of the camera
  from: [x, y, z] # floats, where the camera is located
  to: [x, y, z] # floats, where the camera is aimed at
  up: [x, y, z] # floats, vector starting at the camera and pointing to the cameras up
objects:
- type: sphere | plane | cube | cylinder | cone | triangle | group
  params: # as per the type of the object
          # sphere, plane, cube - no parameters
          # cylinder, cone:  "minimum", "maximum" - floats for cutoff on the Y axis, "closed" - boolean for capping the shape
          # triangle: p1, p2, p3 - [ x, y, z] values for each point of the triangle
          # group: Either:
          #     "objfile" - string pointing to a Wavefront OBJ file location (relative to the CWD)
          #     "content" - exactly the same as the top-level "objects" section
  transform: # optional section, defaults to identity
  - type : identity | translate | scale | rotatex | rotatey | rotatez | shear
    params: # an array of floats that matches the transform type
            # identity - no params
            # translate - [ x, y, z ] floats
            # scale - [ x, y, z ] floats
            # rotate x,y,z - [ radians ] 
            # shear [ xy, xz, yx, yz, zx, zy ] floats
  material:  # Exactly the same as in the above section. Either use a preset, or customize a preset
             # if a name attribute is provided, the resulting material will be saved in the cache
             # potentially overriding any existing cache content
```
