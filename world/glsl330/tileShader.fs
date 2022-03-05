#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in mat4 instance;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform vec2 u_mouse;

// Output fragment color
out vec4 finalColor;

vec4 col;   // rgb dec / 255, a 1
float threshold = 0.025;
float tileSize = 5.0 * 2.0; // tileSize * tileSpacing
vec4 meshPos;
float meshY;

// Change colour depending on Y level
vec4 minCol;
vec4 maxCol;
vec4 diffCol;

void main()
{
    meshPos = instance * vec4(0, 0, 0, 1);
    meshY = meshPos.y;

    minCol = vec4(0.373, 0.388, 0.267, 1);
    maxCol = vec4(0.255, 0.271, 0.184, 1);
    // bigger value - smaller value
    diffCol = vec4(minCol.x - maxCol.x, minCol.y - maxCol.y, minCol.z - maxCol.z, 1);

    // Outline
    if (fragTexCoord.y < threshold || fragTexCoord.y > 1 - threshold || fragTexCoord.x < threshold || fragTexCoord.x > 1 - threshold) {
        col = vec4(0, 0, 0, 1);
    } else {
        if (meshY <= 0.3*tileSize) {
            col = vec4(0, 0.639, 0.8, 1);
        } 
        else {
            col = vec4(minCol.x - (diffCol.x*(meshY/10)), minCol.y - (diffCol.y*(meshY/10)), minCol.z - (diffCol.z*(meshY/10)), 1);
        }

    }
    finalColor = col;
}