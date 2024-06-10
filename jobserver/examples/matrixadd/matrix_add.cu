#include <iostream>
#include <cuda.h>

using std:: cout;

typedef struct{
    double *matriz;
    int    lin;
    int    col;
} Matriz;

__global__ void addMatrix(const Matriz A, const Matriz B, Matriz C)
{
    int idx = threadIdx.x + blockDim.x*blockIdx.x;
    int idy = threadIdx.y + blockDim.y*blockIdx.y;

    if ((idx < A.col) && (idy < A.lin))
      C.matriz[C.col*idy + idx] = A.matriz[A.col*idy + idx] + B.matriz[B.col*idy + idx];
}

void somaMatriz(const Matriz A, const Matriz B, Matriz *C)
{
    Matriz dA;
    Matriz dB;
    Matriz dC;

    int BLOCK_SIZE = 16;

    dA.lin = A.lin;
    dA.col = A.col;
    dB.lin = B.lin;
    dB.col = B.col;
    dC.lin = C->lin;
    dC.col = C->col;

    cudaMalloc((void**)&dA.matriz, dA.lin*dA.col*sizeof(double));
    cudaMalloc((void**)&dB.matriz, dB.lin*dB.col*sizeof(double));
    cudaMalloc((void**)&dC.matriz, dC.lin*dC.col*sizeof(double));

    cudaMemcpy(dA.matriz, A.matriz, dA.lin*dA.col*sizeof(double), cudaMemcpyHostToDevice);
    cudaMemcpy(dB.matriz, B.matriz, dB.lin*dB.col*sizeof(double), cudaMemcpyHostToDevice);

    dim3 dimBlock(BLOCK_SIZE, BLOCK_SIZE);
    dim3 dimGrid((dA.col + dimBlock.x - 1)/dimBlock.x, (dA.lin + dimBlock.y -1)/dimBlock.y);

    addMatrix<<<dimGrid, dimBlock>>>(dA, dB, dC);

    cudaMemcpy(C->matriz, dC.matriz, dC.lin*dC.col*sizeof(double), cudaMemcpyDeviceToHost);
    cudaFree(dA.matriz);
    cudaFree(dB.matriz);
    cudaFree(dC.matriz);

    return;
}

int printMatriz(Matriz *mat){
    for (int y = 0; y < mat->lin; y++)
    {
        for (int x = 0; x < mat->col; x++)
            cout << mat->matriz[y*mat->col + x] << " ";
        cout << "\n";
    }
    return 0;
}

int main(void)
{

    Matriz A;
    Matriz B;
    Matriz *C = new Matriz;
    int lin = 16;
    int col = 7;

    A.lin = lin;
    A.col = col;
    B.lin = lin;
    B.col = col;
    C->lin = lin;
    C->col = col;
    C->matriz = new double[lin*col];

    A.matriz = new double[lin*col];
    B.matriz = new double[lin*col];

    for (int y = 0; y < lin; y++)
        for (int x = 0; x < col; x++)
        {
            A.matriz[y*A.col + x] = 1./(float)(10.*x + y + 10.0);
            B.matriz[y*B.col + x] = (float)(x + y + 1);
        }
    cout << "Matriz A\n";
    printMatriz(&A);
    cout << "Matriz B\n";
    printMatriz(&B);

    somaMatriz(A, B, C);

    cout << "Matriz C\n";
    printMatriz(C);

    return 0;

}