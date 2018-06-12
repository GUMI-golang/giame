#version 430


#define MAXUINT16 65535

layout(r32i, binding = 0) uniform iimage2D dst;
layout (std430, binding = 1) buffer readonly Points {
    float points[];
};
layout (std430, binding = 2) buffer readonly Indexes {
    int indexes[];
};
layout (local_size_x = 1, local_size_y = 1, local_size_z = 1) in;

float startPoint(float f32);
float endPoint(float f32);
float CheckSign(float f32);
int To(float f32);

void main() {
    int start = indexes[int(gl_GlobalInvocationID.y)];
    int end = indexes[int(gl_GlobalInvocationID.y) + 1];
    if (start >= end){
        return;
    }
    float dir = 0;
    float from = 0;
    float to;
    for (int i = start; i < end; i ++){
        to = abs(points[i]);
        if (dir == 0.f){
            dir += CheckSign(points[i]);
            from = to;
            continue;
        }
        float a = startPoint(from);
        float b = endPoint(to);
        if (a == b ){
            imageAtomicExchange(dst, ivec2(int(a), gl_GlobalInvocationID.y), To(to - from));
        }else{
            imageAtomicExchange(dst, ivec2(int(a), gl_GlobalInvocationID.y), To(min(a - from, .5) + .5));
            imageAtomicExchange(dst, ivec2(int(b), gl_GlobalInvocationID.y), To(min(to - b, .5) + .5));
            for(float x = a + 1; x < b; x++) {
                imageAtomicExchange(dst, ivec2(int(x), gl_GlobalInvocationID.y), To(1));
            }
        }
        dir += CheckSign(points[i]);
        from = to;
    }
}


float startPoint(float f32) {
	return float(ceil(f32 - .5f)) + .5f;
}
float endPoint(float f32) {
	return float(floor(f32 + .5f)) - .5f;
}
float CheckSign(float f32) {
    return float(int(floatBitsToUint(f32) >> 31) * -2 + 1);
}
int To(float f32) {
    return clamp(int(f32 * MAXUINT16), 0, MAXUINT16);
}