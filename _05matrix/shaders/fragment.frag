#version 330 core

in vec3 tmpColor;
in vec2 tmpTexCoord;

out vec4 outColor;

uniform sampler2D texture0;
uniform sampler2D texture1;
uniform float rate;

void main()
{
    outColor = mix(texture(texture0, tmpTexCoord), texture(texture1, vec2(tmpTexCoord.x, tmpTexCoord.y)) * vec4(tmpColor, 1.0), rate);
}
