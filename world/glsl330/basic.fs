#version 330

// Input vertex attributes (from vertex shader)
in vec3 fragPosition;
in vec2 fragTexCoord;
in vec3 instance;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform vec2 u_mouse;

// Output fragment color
out vec4 finalColor;

vec4 col;
float threshold = 0.025;
float tileSize = 5.0;
vec3 meshPos;

void main()
{
    meshPos = vec3(instance);

    // Outline
    if (fragTexCoord.y < threshold || fragTexCoord.y > 1 - threshold || fragTexCoord.x < threshold || fragTexCoord.x > 1 - threshold) {
        col = vec4(0, 0, 0, 1);
    } else {

        if (meshPos.y <= 0.3) {
            col = vec4(1, 0, 0, 1);
        } else if (meshPos.y < 0.6) {
            col = vec4(0, 1, 0, 1);
        } else {
            col = vec4(0, 0, 1, 1);
        }

        // col = vec4(model_matrix, 1);
    }
    finalColor = col;
}
