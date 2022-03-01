#version 330

// Input vertex attributes (from vertex shader)
in vec3 fragPosition;
in vec2 fragTexCoord;
in mat4 instance;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform vec2 u_mouse;

uniform mat4 view_matrix;
uniform mat4 model_matrix;


// Output fragment color
out vec4 finalColor;

vec4 col;
float threshold = 0.025;
float tileSize = 5.0;
vec3 meshPos;

void main()
{
    meshPos = vec3(instance * model_matrix);

    // Outline
    if (fragTexCoord.y < threshold || fragTexCoord.y > 1 - threshold || fragTexCoord.x < threshold || fragTexCoord.x > 1 - threshold) {
        col = vec4(0, 0, 0, 1);
    } else {

        if (meshPos[2] <= 0.3) {
            col = vec4(1, 0, 0, 1);
        } else if (meshPos[2] < 0.6) {
            col = vec4(0, 1, 0, 1);
        } else {
            col = vec4(0, 0, 1, 1);
        }

        // col = vec4(model_matrix, 1);
    }
    finalColor = col;
}
