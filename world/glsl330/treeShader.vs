#version 330

// Input vertex attributes
in vec3 vertexPosition;
in vec2 vertexTexCoord;

in mat4 instanceTransform;

// Input uniform values
uniform mat4 mvp;

// Output vertex attributes (to fragment shader)
out vec2 fragTexCoord;

// NOTE: Add here your custom variables 

void main()
{
    // Compute MVP for current instance
    mat4 mvpi = mvp*instanceTransform;
    
    // Send vertex attributes to fragment shader
    fragTexCoord = vertexTexCoord;

    // Calculate final vertex position
    gl_Position = mvpi*vec4(vertexPosition, 1);
}