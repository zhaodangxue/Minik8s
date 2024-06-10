#!/bin/bash

#SBATCH --job-name=test
#SBATCH --partition=dgx2
#SBATCH -N 1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=6
#SBATCH --gres=gpu:1              #若使用 2 块卡，就给 gres=gpu:2
#SBATCH --mail-type=end
#SBATCH --mail-user=varia07@sjtu.edu.cn
#SBATCH --output=output/%j.out
#SBATCH --error=output/%j.err

module load gcc cuda

nvcc -o matrix_add matrix_add.cu
./matrix_add
