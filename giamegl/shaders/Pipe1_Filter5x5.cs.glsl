#version 430

#define MAXUINT16 65535

layout(r32i, binding = 0) uniform coherent iimage2D to;
layout(std430, binding = 1) buffer _Mask{
    float [5][5]mask;
};

layout (local_size_x = 1, local_size_y = 1) in;


void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
    float c;
    // y = -2
    c += float(imageLoad(to, pos + ivec2(-2, -2)).x) * mask[0][0];
    c += float(imageLoad(to, pos + ivec2(-1, -2)).x) * mask[1][0];
    c += float(imageLoad(to, pos + ivec2(+0, -2)).x) * mask[2][0];
    c += float(imageLoad(to, pos + ivec2(+1, -2)).x) * mask[3][0];
    c += float(imageLoad(to, pos + ivec2(+2, -2)).x) * mask[4][0];
    // y = -1
    c += float(imageLoad(to, pos + ivec2(-2, -1)).x) * mask[0][1];
    c += float(imageLoad(to, pos + ivec2(-1, -1)).x) * mask[1][1];
    c += float(imageLoad(to, pos + ivec2(+0, -1)).x) * mask[2][1];
    c += float(imageLoad(to, pos + ivec2(+1, -1)).x) * mask[3][1];
    c += float(imageLoad(to, pos + ivec2(+2, -1)).x) * mask[4][1];
    // y = 0
    c += float(imageLoad(to, pos + ivec2(-2, +0)).x) * mask[0][2];
    c += float(imageLoad(to, pos + ivec2(-1, +0)).x) * mask[1][2];
    c += float(imageLoad(to, pos + ivec2(+0, +0)).x) * mask[2][2];
    c += float(imageLoad(to, pos + ivec2(+1, +0)).x) * mask[3][2];
    c += float(imageLoad(to, pos + ivec2(+2, +0)).x) * mask[4][2];
    // y = +1
    c += float(imageLoad(to, pos + ivec2(-2, +1)).x) * mask[0][3];
    c += float(imageLoad(to, pos + ivec2(-1, +1)).x) * mask[1][3];
    c += float(imageLoad(to, pos + ivec2(+0, +1)).x) * mask[2][3];
    c += float(imageLoad(to, pos + ivec2(+1, +1)).x) * mask[3][3];
    c += float(imageLoad(to, pos + ivec2(+2, +1)).x) * mask[4][3];
    // y = +2
    c += float(imageLoad(to, pos + ivec2(-2, +2)).x) * mask[0][4];
    c += float(imageLoad(to, pos + ivec2(-1, +2)).x) * mask[1][4];
    c += float(imageLoad(to, pos + ivec2(+0, +2)).x) * mask[2][4];
    c += float(imageLoad(to, pos + ivec2(+1, +2)).x) * mask[3][4];
    c += float(imageLoad(to, pos + ivec2(+2, +2)).x) * mask[4][4];
    imageStore(to, pos, ivec4(clamp(c, 0, MAXUINT16), 0,0,0));
}
