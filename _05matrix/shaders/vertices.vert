#version 330 core

layout(location = 0) in vec3 inPosition;
layout(location = 1) in vec3 inColor;
layout(location = 2) in vec2 inTexCoord;

out vec3 tmpColor;
out vec2 tmpTexCoord;

uniform mat4 tran;

void main()
{
    gl_Position = tran * vec4(inPosition, 1.0);
    tmpColor = inColor;
    tmpTexCoord = inTexCoord;
}
