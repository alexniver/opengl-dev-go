#version 330 core

in vec3 tmpColor;
in vec2 tmpTexCoord;

out vec4 outColor;

uniform sampler2D texture1;

void main()
{
    outColor = texture(texture1, tmpTexCoord) * vec4(tmpColor, 1.0);
}
