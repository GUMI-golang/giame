#version 430

#define MAXUINT16 65535

layout(rgba32f, binding = 0) uniform image2D result;
layout(r32i, binding = 1) uniform coherent iimage2D mask;
layout(rgba32f, binding = 2) uniform image2D filler;
layout(std430, binding = 3) buffer RepeatInfo{
    ivec2 start;
    ivec2 repeator;
};
layout (local_size_x = 1, local_size_y = 1) in;

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
    ivec2 fillerpos = pos;
    ivec2 fillersize = imageSize(filler);
    fillerpos.x %= repeator.x;
    fillerpos.y %= repeator.y;
    float intense = float(imageLoad(mask, pos).x) / MAXUINT16;


    ivec2 resultsize = imageSize(result);
    ivec2 resultpos = ivec2(pos + start);
    resultpos.y = resultsize.y - resultpos.y - 1;

    vec4 prev = imageLoad(result, resultpos);
    vec4 need = imageLoad(filler, fillerpos) * intense;
    imageStore(result, resultpos, mix(prev, need, need.w));
}
