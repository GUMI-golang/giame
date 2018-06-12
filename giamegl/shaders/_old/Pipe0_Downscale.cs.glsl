#version 430

#define MAXUINT16 65535

layout(r32i, binding = 0) uniform coherent iimage2D to;
layout(r32i, binding = 1) uniform coherent iimage2D from;

layout (local_size_x = 1, local_size_y = 1) in;

int wsg(ivec2 toSize, ivec2 fromSize, ivec2 toPos);
float fn(float x);

void main() {
    ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
    imageStore(to, pos, ivec4(wsg(imageSize(to), imageSize(from), pos), 0,0,0));
}

float fn(float x){
    x = abs(x);
	if (x <= 1.) {
		return 2*x*x*x - 3*x*x + 1;
	}
	return 0;
}

int wsg(ivec2 toSize, ivec2 fromSize, ivec2 toPos){
    vec2 delta = vec2(
        float(fromSize.x) / float(toSize.x),
        float(fromSize.y) / float(toSize.y)
    );
    vec2 scale = vec2(
        float(toSize.x) / float(fromSize.x),
        float(toSize.y) / float(fromSize.y)
    );
    vec2 halfscale = scale/2;
    ivec2 a = ivec2(int(delta.x), int(delta.y));
    //
    float temp;
    float sum;
    vec2 s = (toPos - halfscale);
    s.x = s.x * delta.x;
    s.y = s.y * delta.y;
    for (int i = -a.x; i <= a.x; i++ ){
        for (int j = -a.y; j <= a.y; j++ ){
            vec2 t = s + vec2(float(i), float(j)) ;
            ivec2 it = ivec2(t);
            if(it.x <0 || it.x >= fromSize.x){
                continue;
            }
            if(it.y <0 || it.y >= fromSize.y){
                continue;
            }
            float kr = fn(float(it.x) - s.x - 0.00001) * fn(float(it.y) - s.y - 0.00001);

            if(-0.0001 < kr && kr < 0.0001){
                continue;
            }
            temp += float(imageLoad(from, it).x) * kr;
            sum += kr;
        }
    }
    return int(temp / sum);
}
