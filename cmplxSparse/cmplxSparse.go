// Package for sparse complex matrix operations for NEGF calculations

package cmplxSparse

import (
	"errors"
	"math"
	"math/cmplx"
	"fmt"
)

// Sparse matrix data structure stores a matrix in Diagonal form, which is convenient
// since the Hamiltonian is built from a tridiagonal matrix. Each column of the 2D-array
// given by Data is a diagonal of the actual matrix. Hence, the number of columns must be
// odd. The middle column is the main diagonal. Every column right of the middle is the
// off-diagonal to the right, while every column left of the middle is the off-diagonal
// to the left. The row index of an entry corresponds to the matrix row the entry belongs.
// i.e. the identity matrix is given as:
//      A[i][0] = 1.0
// i.e. the standard 4x4 tridiagonal matrix is given as:
//      A =  0.0 2.0 -1.0
//          -1.0 2.0 -1.0
//          -1.0 2.0 -1.0
//          -1.0 2.0  0.0

type sparseMat struct {
	Data			[][]complex128
}

func init() {
}

func New() *sparseMat {
	tmp := new(sparseMat);
	return tmp;
}

// Function to access matrix elements
func AccessMatrix( m, n int, s *sparseMat ) complex128 {
	if ((m < 0) || (n < 0)) {
		errors.New("Access indices cannot be negative!");
	} else {
		matDim := len(s.Data)
		if ((m >= matDim) || (n >= matDim)) {
			errors.New("Access index out of range!");
		}
	}
	totalDiags := len(s.Data[0]);
	maxOffset := (totalDiags-1)/2;
	if ((n > m+maxOffset) || (n < m-maxOffset)) {
		return 0.0
	} else {
		idx0 := n-m+maxOffset
		return s.Data[m][idx0]
	}
}

// Function to print matrix to screen
func PrintSparseMatrix( s *sparseMat ) {
	matSize := len(s.Data);
	if (matSize < 1) {
		return;
	}
        matDiags := len(s.Data[0]);
	fmt.Println("----------------------------------------");
	if (matDiags == 1) {
		fmt.Println("The matrix is 1 x 1 :");
		fmt.Println(s.Data[0][0]);
	} else {
		fmt.Println("The matrix is", matSize, "x", matSize, ":");
		for idx0 := 0; idx0 < matSize; idx0++ {
			for idx1 := 0; idx1 < matSize; idx1++ {
				dataValue := AccessMatrix(idx0, idx1, s);
				fmt.Printf("%f    ",dataValue);
			}
			fmt.Printf("\n");
		}
	}	
	fmt.Println("----------------------------------------");
}

// Function to initialize an identity matrix in Diagonal format
func MakeIdentity( matSize int, s *sparseMat ) {
	// TODO: error if matsize is zero and less
	if (matSize <= 0) {
		errors.New("ERROR: Invalid size for matrix (<=0)");
	}
	s.Data = make([][]complex128, matSize);
	for idx0 := 0; idx0 < matSize; idx0++ {
		s.Data[idx0] = make([]complex128,1);
		s.Data[idx0][0] = 1.0;
	}
}

// Function to initialize a standard tridiagonal matrix
func MakeTriDiag ( matSize int, s *sparseMat ) {
	// TODO: error if matSize is zero and less
	if (matSize <= 0) {
		errors.New("ERROR: matrix size must be 1 or larger!");
	}
	s.Data = make([][]complex128, matSize);
	for idx0 := 0; idx0 < matSize; idx0++ {
		if (matSize == 1) {
			s.Data[idx0] = make([]complex128,1);
			s.Data[idx0][0] = 2.0;
		} else {
			s.Data[idx0] = make([]complex128,3);
			if (idx0 == 0) {
				s.Data[idx0][0],s.Data[idx0][1],s.Data[idx0][2] =  0.0, 2.0, -1.0;
			} else if (idx0 == matSize - 1) {
				s.Data[idx0][0],s.Data[idx0][1],s.Data[idx0][2] = -1.0, 2.0,  0.0;
			} else {
				s.Data[idx0][0],s.Data[idx0][1],s.Data[idx0][2] = -1.0, 2.0, -1.0;
			}
		}
	}
}

// Function for creating basic tridiagonal Hamiltonian for MTJ
func MakeHamTriDiag( grdSize int, s *sparseMat ) {
	if (grdSize < 5) {
		errors.New("ERROR: Creating the Hamiltonian requires at least 5 points!");
	}
	s.Data = make([][]complex128, 2*grdSize);
	for idx0 := 0; idx0 < 2*grdSize; idx0++ {
		s.Data[idx0] = make([]complex128,5);
		s.Data[idx0][0], s.Data[idx0][1], s.Data[idx0][2], s.Data[idx0][3], s.Data[idx0][4] = -1.0, 0.0, 2.0, 0.0, -1.0;
	}
	for idx0 := 0; idx0 < 2; idx0++ {
		offSet := 2*grdSize-1-idx0;
		s.Data[idx0][0], s.Data[offSet][4] = 0.0, 0.0;
	}
}

// Function for adding two sparse matrices stored in diagonal format
func SparseDiagAdd(s, t *sparseMat) *sparseMat {
	Ssize, Tsize := len(s.Data), len(t.Data);
	if (Ssize != Tsize) {
		errors.New("ERROR: cannot add matrices of different sizes!");
	}

	mainDiagIdxS, mainDiagIdxT := (len(s.Data[0]) - 1)/2, (len(t.Data[0]) - 1)/2;

	u := t;
	flag0 := 1;
	sweepIdx := mainDiagIdxS;

	if (mainDiagIdxT < mainDiagIdxS) {
		u = s;
		flag0 = 0;
		sweepIdx = mainDiagIdxT;
	}

	if (flag0 == 0) {
		// u is storing s
	} else {
		// u is storing t
	}
	for idx0 := 0; idx0 < Ssize; idx0++ {
		// Calculate the main diagonal first
		if (flag0 == 0) {
			// u is storing s
			u.Data[idx0][mainDiagIdxS] += t.Data[idx0][mainDiagIdxT];
		} else {
			// u is storing t
			u.Data[idx0][mainDiagIdxT] += t.Data[idx0][mainDiagIdxS];
		}
		// Calculate the off-diagonal terms
		termIdx := idx0;
		if (sweepIdx < termIdx) {
			termIdx = sweepIdx;
		}
		for idx1 := 1; idx1 < termIdx; idx1++ {
			if (flag0 == 0) {
				// u is storing s
				u.Data[idx0][mainDiagIdxS - idx1] += t.Data[idx0][mainDiagIdxT - idx1];
			} else {
				// u is storing t
				u.Data[idx0][mainDiagIdxT - idx1] += t.Data[idx0][mainDiagIdxS - idx1];
			}
		}
	}

	return u;
}

// Function for scaled sparse matrix (t) to another sparse matrix (s), both stored in diagonal format
func SparseDiagMAdd(mult complex128, s, t *sparseMat) *sparseMat {
	Ssize, Tsize := len(s.Data), len(t.Data);
	if (Ssize != Tsize) {
		errors.New("ERROR: cannot add matrices of different sizes!");
	}

	u := t;
	ScaleSparseMatrixIP(mult, u);
	u = SparseDiagAdd(s,u);
	return u;
}

// Function for scalar multiplication of sparse matrix
func ScaleSparseMatrix(A complex128, B *sparseMat) *sparseMat {
	s := B;

	for idx0 := 0 ; idx0 < len(s.Data); idx0++ {
		for idx1 := 0; idx1 < len(s.Data[idx0]); idx1++ {
			s.Data[idx0][idx1] *= A;
		}
	}

	return s
}

// Function for in-place scalar multiplication of sparse matrix
func ScaleSparseMatrixIP(A complex128, s *sparseMat) {
	for idx0 := 0 ; idx0 < len(s.Data); idx0++ {
		for idx1 := 0; idx1 < len(s.Data[idx0]); idx1++ {
			s.Data[idx0][idx1] *= A;
		}
	}
}

// Function for adding applied voltage profile to Hamiltonian
func AddVoltagePotential( N_fm, N_ox int, voltage float64, s *sparseMat ) *sparseMat {
	tmpLength := len(s.Data);
	totalPts := tmpLength/2;
	voltageProfile := make([]float64,totalPts)
	voltageDelta := voltage/float64(N_ox+1);
	t := s;
	for idx0 := 0; idx0 < totalPts; idx0++ {
		if (idx0 <= N_fm) {
			voltageProfile[idx0] = 0.5*voltage;
		} else if (idx0 < N_fm+1+N_ox) {
			voltageProfile[idx0] = voltageProfile[idx0-1]-voltageDelta;
		} else {
			voltageProfile[idx0] = -0.5*voltage;
		}
		t.Data[2*idx0][2] += complex(voltageProfile[idx0], 0.0);
		t.Data[2*idx0+1][2] += complex(voltageProfile[idx0], 0.0);
	}
	return t
}

// Function for adding barrier potential profile to Hamiltonian
func AddBarrierProfile( N_fm, N_ox int, Eb float64, s *sparseMat ) {
	tmpLength := len(s.Data);
	totalPts := tmpLength/2;
	if ((N_fm >= totalPts) || (N_ox >= totalPts)) {
		errors.New("ERROR: Indices are out of range!");
	} else if ((N_fm < 0) || (N_ox < 0)) {
		errors.New("ERROR: Indices cannot be negative!");
	}
	s.Data[2*N_fm][2] += complex(0.5*Eb,0);
	s.Data[2*N_fm+1][2] += complex(0.5*Eb,0);
	s.Data[2*(N_fm+N_ox+1)][2] += complex(0.5*Eb,0);
	s.Data[2*(N_fm+N_ox)+3][2] += complex(0.5*Eb,0);
	for idx0 := N_fm+1; idx0 < N_fm+1+N_ox; idx0++ {
		s.Data[2*idx0][2] += complex(Eb,0);
		s.Data[2*idx0+1][2] += complex(Eb,0);
	}
}

// Function for adding band splitting for up-spin and down-spin conduction bands on left contact
func AddBandSplitLeftFM(mx, my, mz, deltE float64, N_fm int, s *sparseMat) *sparseMat {
	Mnorm := mx*mx + my*my + mz*mz;
	if (Mnorm == 0) {
		errors.New("ERROR: invalid mx, my and mz combination!");
	} else if (N_fm < 0) {
		errors.New("ERROR: invalid value of N_fm!");
	}
	Mnorm = math.Sqrt(Mnorm);
	m_x := mx/Mnorm;
	m_y := my/Mnorm;
	m_z := mz/Mnorm;
	BT := make([][]complex128,2);
	for idx0 := 0; idx0 < 2; idx0++ {
		BT[idx0] = make([]complex128,2);
	}
	scaleFac := 0.5*deltE;
	BT[0][0] = complex(scaleFac*(1.0 - m_z),0.0);
	BT[1][1] = complex(scaleFac*(1.0 + m_z),0.0);
	BT[0][1] = complex(m_x*scaleFac,-m_y*scaleFac);
	BT[1][0] = complex(m_x*scaleFac,m_y*scaleFac);

	t := s;

	for idx0 := 0; idx0 < N_fm; idx0++ {
		currIdx := 2*idx0;
		t.Data[currIdx][2] += BT[0][0];
		t.Data[currIdx+1][2] += BT[1][1];
		t.Data[currIdx][3] += BT[0][1];
		t.Data[currIdx+1][1] += BT[1][0];
	}
	t.Data[2*N_fm][2] += 0.5*BT[0][0];
	t.Data[2*N_fm+1][2] += 0.5*BT[1][1];
	t.Data[2*N_fm][3] += 0.5*BT[0][1];
	t.Data[2*N_fm+1][1] += 0.5*BT[1][0];

	return t;
}

// Function for adding band splitting for up-spin and down-spin conduction bands on right contact
func AddBandSplitRightFM(mx, my, mz, deltE float64, N_fm int, s *sparseMat) *sparseMat {
	Mnorm := mx*mx + my*my + mz*mz;
	if (Mnorm == 0) {
		errors.New("ERROR: invalid mx, my and mz combination!");
	} else if (N_fm < 0) {
		errors.New("ERROR: invalid value of N_fm!");
	}
	Mnorm = math.Sqrt(Mnorm);
	m_x := mx/Mnorm;
	m_y := my/Mnorm;
	m_z := mz/Mnorm;
	BT := make([][]complex128,2);
	for idx0 := 0; idx0 < 2; idx0++ {
		BT[idx0] = make([]complex128,2);
	}
	scaleFac := 0.5*deltE;
	BT[0][0] = complex(scaleFac*(1.0 - m_z),0.0);
	BT[1][1] = complex(scaleFac*(1.0 + m_z),0.0);
	BT[0][1] = complex(m_x*scaleFac,-m_y*scaleFac);
	BT[1][0] = complex(m_x*scaleFac,m_y*scaleFac);

	t := s;

	grdSz := len(s.Data)/2;
	startIdx := grdSz-N_fm;
	for idx0 := startIdx; idx0 < grdSz; idx0++ {
		currIdx := 2*idx0;
		t.Data[currIdx][2] += BT[0][0];
		t.Data[currIdx+1][2] += BT[1][1];
		t.Data[currIdx][3] += BT[0][1];
		t.Data[currIdx+1][1] += BT[1][0];
	}
	t.Data[2*(startIdx-1)][2] += 0.5*BT[0][0];
	t.Data[2*startIdx-1][2] += 0.5*BT[1][1];
	t.Data[2*(startIdx-1)][3] += 0.5*BT[0][1];
	t.Data[2*startIdx-1][1] += 0.5*BT[1][0];

	return t;
}

// Function for in-place scalar multiplication of part of sparse Hamiltonian matrix
// for placing t_0 in correct positions
func ScaleRangeSparseMatrixIP(startPtr, endPtr int, diagIdx int, A complex128, s *sparseMat) {
	grdPts := len(s.Data);
	diagNum := (len(s.Data[0]) + 1)/2;

	if ((startPtr < 0) || (endPtr < 0)) {
		errors.New("ERROR: startPtr and endPtr cannot be negative!");
	} else if (startPtr > endPtr) {
		errors.New("ERROR: startPtr cannot be greater than endPtr!");
	} else if ((startPtr >= grdPts) || (endPtr >= grdPts)) {
		errors.New("ERROR: startPtr and endPtr out of range!");
	}

	if (int(math.Abs(float64(diagIdx))) < diagNum) {
		targIdx := diagNum - 1 + diagIdx;
		if (startPtr == endPtr) {
			s.Data[startPtr][targIdx] *= A;
		} else {
			for idx0 := startPtr ; idx0 <= endPtr; idx0++ {
				s.Data[idx0][targIdx] *= A;
			}
		}
	}
}

// Special function for accessing elements of sparse matrices
// stored in Diagonal format
func SparseDiagAccess( m, n int, s *sparseMat ) complex128 {
	maxSize := len(s.Data);
	if ((m >= maxSize) || (n >= maxSize)) {
		errors.New("ERROR: Indices are out of range!");
	} else if ((m < 0) || (n < 0)) {
		errors.New("ERROR: Indices are negative!");
	}
	maxIdx := len(s.Data[0]) - 1
	baseIdx := maxIdx/2;
	offSet := n-m;
	colIdx := baseIdx+offSet;
	if ((colIdx < 0) || (colIdx > maxIdx)) {
		return complex(0.0, 0.0);
	} else {
		return s.Data[m][colIdx]
	}
}

// Function for executing Doolittle's algorithm to perform LU
// factorization of sparse matrix stored in Diagonal format
func SparseDiagLU(s *sparseMat) *sparseMat {
	// Get indices for performing for loops...
	maxSize := len(s.Data);
	t := s;

	// Handle small matrix sizes (up to 2 x 2)
	if (maxSize == 1) {
		return t;
	} else if (maxSize == 2) {
		if (s.Data[1][0] == complex(0.0,0.0)) {
			return t;
		} else {
			t.Data[1][0] /= t.Data[1][0];
			t.Data[1][1] -= t.Data[0][2]*t.Data[0][1]/t.Data[1][0];
			return t; 
		}
	}
	maxMatIdx := len(s.Data[0]) - 1
	mainDiagIdx := maxMatIdx/2;

	// For matrices 3x3 and larger

	var (
		endIdx		int
		topDiagIdx	int
		t_num_R		float64
		t_num_I		float64
		t_num		float64
		targRow		int
		targCol		int
	)

	// The following has been optimized for the assumption that s.Data stores equal number of diagonals
	// above and below the main diagonal. i.e., s.Data must have odd number of columns. The middle
	// column are the entries for main diagonal. The n-th column to the left and right of the middle
	// column are the n-th diagonal below and above the main diagonal, respectively.

	// Scan down main diagonal...
	for diagIdx := 1; diagIdx < maxSize; diagIdx++ {

		// First, calculate the column of the L matrix...

		// Determine number of elements down the current column of the L matrix
		// that needs to be calculated
		endIdx = maxSize - diagIdx;
		if (mainDiagIdx < endIdx) {
			endIdx = mainDiagIdx;
		}
		topDiagIdx = diagIdx-1;
		t_num_R, t_num_I = real(t.Data[topDiagIdx][mainDiagIdx]), imag(t.Data[topDiagIdx][mainDiagIdx]);
		t_num = t_num_R*t_num_R + t_num_I*t_num_I;

		for LColIdx := 1; LColIdx <= endIdx; LColIdx++ {
			targRow = topDiagIdx + LColIdx;
			targCol = mainDiagIdx - LColIdx;

			// Only scan leftwards of L row if there are non-zero entries to the left of
			// current entry of L matrix
			if ((diagIdx > 1) && (targCol > 0)) {
				for colIdx := 1; colIdx < LColIdx; colIdx++ {
					t.Data[targRow][targCol] -= t.Data[targRow][targCol-colIdx]*t.Data[targRow-1-colIdx][targCol+1+colIdx];
				}
			}

			// After subtracting, divide by the corresponding element on main diagonal of
			// U matrix
			t.Data[targRow][targCol] *= complex(t_num_R/t_num, -1.0*t_num_I/t_num);
		}

		// After calculating the column of L matrix, we can move on to
		// calculate the row of the U matrix...

		for currColIdx := mainDiagIdx; currColIdx < mainDiagIdx + endIdx; currColIdx++ {
			totalRight := maxMatIdx - currColIdx;
			RowScanEndIdx := diagIdx;
			if (totalRight < RowScanEndIdx) {
				RowScanEndIdx = totalRight;
			}
			if (RowScanEndIdx > 0) {
				for scanIdx := 1; scanIdx <= RowScanEndIdx; scanIdx++ {
					t.Data[diagIdx][currColIdx] -= t.Data[diagIdx-scanIdx][currColIdx+scanIdx]*t.Data[diagIdx][mainDiagIdx-scanIdx];
				}
			}
		}
	}

	return t;
}

// Function to solvee A*x = b, where sparse matrix stored in diagonal format
func SparseDiagLinearSolver(A *sparseMat, b []complex128) []complex128 {
	MatrixSize, VectorSize := len(A.Data), len(b);
	if (MatrixSize != VectorSize) {
		errors.New("ERROR: Mismatch between matrix size and vector length!");
	}

	// Find LU factorization of A first
	LU_A := SparseDiagLU(A);

	// Use back substitution to determine x:
	// First solve L*y = b using back substitution. Then, solve U*x = y.

	// Create variables for indexing
	numMatCols := len(A.Data[0]);
	mainDiagIdx := (numMatCols - 1)/2;
	var (
		idx0, loopCnt, offSetIdx		int
	)

	// Use buffer to store result
	x := b;

	for idx0 = 1; idx0 < MatrixSize; idx0++ {
		loopCnt = idx0;
		if (mainDiagIdx < idx0) {
			loopCnt = mainDiagIdx;
		}

		for offSetIdx = 1; offSetIdx < loopCnt+1; offSetIdx++ {
			x[idx0] -= x[(idx0-offSetIdx)]*LU_A.Data[idx0][(mainDiagIdx-offSetIdx)];
		}
	}

	for idx0 = 0; idx0 < MatrixSize; idx0++ {
		if (idx0 > 1) {
			loopCnt = idx0;
			if (mainDiagIdx < loopCnt) {
				loopCnt = mainDiagIdx;
			}
			for offSetIdx = 1; offSetIdx < loopCnt+1; offSetIdx++ {
				x[idx0] -= x[(idx0+offSetIdx)] * LU_A.Data[idx0][(mainDiagIdx+offSetIdx)];
			}
			t_num_R, t_num_I := real(LU_A.Data[idx0][mainDiagIdx]), imag(LU_A.Data[idx0][mainDiagIdx]);
			t_num := t_num_R*t_num_R + t_num_I*t_num_I;
			x[idx0] *= complex(t_num_R/t_num, -1.0*t_num_I/t_num); 
		}
	}
	return x;
}

// Function to compute basis transformation matrix
func BasisTransform(th, phi float64) *[2][2]complex128 {

	BTMatrix := new([2][2]complex128);
	th_2, phi_2 := 0.5*th, 0.5*phi;
	csn := complex(math.Cos(th_2), 0.0) * cmplx.Exp(complex(0.0, -1.0*phi_2));
	sn  := complex(math.Sin(th_2), 0.0) * cmplx.Exp(complex(0.0, phi_2));

	BTMatrix[0][0] = cmplx.Conj(csn);
	BTMatrix[0][1] = cmplx.Conj(sn);
	BTMatrix[1][0] = complex(-1.0, 0.0)*sn;
	BTMatrix[1][1] = csn;

	return BTMatrix;
}

// Function to obtain complex conjugate (dagger operation) of
// sparse matrix stored in diagonal form
func SparseDiagDagger(s *sparseMat) *sparseMat {

	// Initiaize variables for indexing purposes
	matSize := len(s.Data);
	numDiags := len(s.Data[0]);
	mainDiagIdx := (numDiags-1)/2;

	// Initialize output variable
	t := s;

	for idx0 := 0; idx0 < matSize; idx0++ {
		// Complex conjugate main diagonal
		t.Data[idx0][mainDiagIdx] = cmplx.Conj(t.Data[idx0][mainDiagIdx]);

		// When there are more than one diagonal stored, and if we are not in the top left most entry,
		// we will need to swap rows and columns, and complex conjugate those entries as well.

		if ((idx0 > 1) && (mainDiagIdx > 0)) {
			loopIdx := mainDiagIdx;
			if (idx0 < loopIdx) {
				loopIdx = idx0;
			}
			for idx1 := 1; idx1 < loopIdx; idx1++ {
				// Swap row and columns, and conjugate those entries
				t.Data[idx0][mainDiagIdx-idx1], t.Data[idx0-idx1][mainDiagIdx+idx1] = cmplx.Conj(t.Data[idx0-idx1][mainDiagIdx+idx1]), cmplx.Conj(t.Data[idx0][mainDiagIdx-idx1]);
			}
		}
	}

	return t;
}

