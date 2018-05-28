#version 430

#define MAXUINT16 65535


layout(rgba32f, binding = 0) uniform image2D result;
layout(r32i, binding = 1) uniform coherent iimage2D workspace;
// workspace
layout(std430, binding = 2) buffer readonly _start {
    ivec2 start;
};
layout(std430, binding = 3) buffer readonly _color {
    vec4 color;
};
layout (local_size_x = 1, local_size_y = 1) in;

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
    ivec2 resultpos = pos + start;
    float intense = float(imageLoad(workspace, pos).x) / MAXUINT16;
    // flip vertical
    resultpos.y = imageSize(result).y - int(resultpos.y) - 1;
    //
    vec4 prev = imageLoad(result, resultpos);
    vec4 need = color * intense;
    imageStore(result, resultpos, mix(prev, need, need.w));
}