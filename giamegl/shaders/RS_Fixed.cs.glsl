#version 430

#define MAXUINT16 65535

layout(rgba32f, binding = 0) uniform image2D result;
layout(r32i, binding = 1) uniform coherent readonly iimage2D mask;
layout(rgba32f, binding = 2) uniform readonly image2D filler;
layout(std430, binding = 3) buffer readonly _start{
    ivec2 start;
};
layout (local_size_x = 1, local_size_y = 1) in;

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
    ivec2 fillerpos = pos;
    float intense = float(imageLoad(mask, pos).x) / MAXUINT16;
    ivec2 resultpos = ivec2(pos + start);
    resultpos.y = imageSize(result).y - resultpos.y - 1;

    vec4 prev = imageLoad(result, resultpos);
    vec4 need = imageLoad(filler, fillerpos) * intense;
    imageStore(result, resultpos, mix(prev, need, need.w));
}
