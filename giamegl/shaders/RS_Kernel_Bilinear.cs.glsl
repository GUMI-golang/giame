#version 430

#define MAXUINT16 65535
#define A 1

layout(rgba32f, binding = 0) uniform image2D result;
layout(r32i, binding = 1) uniform coherent iimage2D mask;
layout(rgba32f, binding = 2) uniform image2D filler;
layout(std430, binding = 3) buffer StartDelta{
    ivec2 start;
    vec2 delta;
};
layout (local_size_x = 1, local_size_y = 1) in;
float fn(float x);

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
    //
    ivec2 srcsz = imageSize(filler);
    vec2 spos = vec2(
        delta.x * float(pos.x) - .5,
        delta.y * float(pos.y) - .5
    );
    int[2] hori = {int(spos.x - A), int(spos.x + A + .5)};
    int[2] vert = {int(spos.y - A), int(spos.y + A + .5)};
    //
    vec4 color;
    float sum;
    for(int tx = hori[0]; tx <= hori[1]; tx++) {
        for(int ty = vert[0]; ty <= vert[1]; ty++) {
            int rx = clamp(tx, 0, srcsz.x-1);
            int ry = clamp(ty, 0, srcsz.y-1);

            float kr = fn(float(tx)-spos.x) * fn(float(ty)-spos.y);
            color += imageLoad(filler, ivec2(rx, ry)) * kr;
            sum += kr;
        }
    }
    color = color / sum;
    //
    float mask = clamp(float(imageLoad(mask, pos).x) / MAXUINT16, 0, 1);
    color = color * mask;
    //
    ivec2 resultpos = pos + start;
    resultpos.y = imageSize(result).y - int(resultpos.y) - 1;
    imageStore(result, resultpos, mix(imageLoad(result, pos + start), color, color.w));
}

float fn(float x){
	x = abs(x);
	if (x < 1){
		return 1 - x;
	}
	return 0;
}