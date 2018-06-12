#version 430

#define CLOSETOZERO 0.00001
#define MAXUINT16 65535

layout(r32i, binding = 0) uniform iimage2D ioutput;
layout (std430, binding = 1) buffer readonly Points {
    vec2 points[];
};
layout (local_size_x = 1, local_size_y = 1) in;

void doSimple(vec2 a, vec2 b, float dir);
float areaRemain(vec2 a, vec2 b);
float area(vec2 a, vec2 b);
ivec2 location(vec2 a);
int Percent(float a);

void main() {
    vec2 from, to, tempv2;
    from = points[gl_GlobalInvocationID.x];
    to =  points[(gl_GlobalInvocationID.x + 1)];
    if (isnan(from.x )|| isnan(to.x)){
        return;
    }
    float dir = 1;
    if (to.y < from.y){
        dir = -1;
        tempv2 = from;
        from = to;
        to = tempv2;
    }
    if(to.y - from.y < CLOSETOZERO){
        return;
    }
    //
    float delta = (to.x - from.x) / (to.y - from.y);
    vec2 a, b;
    a = from;
    b = vec2(
        a.x + delta * (ceil(a.y) - a.y),
        ceil(a.y)
    );
    for(; b.y <= to.y; b = vec2(a.x + delta, a.y + 1)){
        doSimple(a, b, dir);
        a = b;
    }
    b = to;
    doSimple(a, b, dir);
}
void doSimple(vec2 a, vec2 b, float dir){
    if((b.y - a.y) < CLOSETOZERO){
        return;
    }
}

float areaRemain(vec2 a, vec2 b){
    float xceil = ceil(max(a.x, b.x));
    float dy = abs(b.y - a.y);
    return dy - ((xceil - a.x) + (xceil - b.x)) * dy / 2;
}

float area(vec2 a, vec2 b){
    float xceil = ceil(max(a.x, b.x));
    return ((xceil - a.x) + (xceil - b.x)) * abs(b.y - a.y) / 2;
}

int Percent(float a){
    return clamp(int(a * MAXUINT16), -MAXUINT16, MAXUINT16);
}

ivec2 location(vec2 a){
    ivec2 temp = ivec2(a);
    return clamp(temp, ivec2(0, 0), imageSize(ioutput));
}