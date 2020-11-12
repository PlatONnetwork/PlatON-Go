// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package math

import (
	"errors"
	"fmt"
	"math"
)

var (
	DELTA         = [...]float64{0.08333333333333333, -2.777777777777778e-5, 7.936507936507937e-8, -5.952380952380953e-10, 8.417508417508329e-12, -1.917526917518546e-13, 6.410256405103255e-15, -2.955065141253382e-16, 1.7964371635940225e-17, -1.3922896466162779e-18, 1.338028550140209e-19, -1.542460098679661e-20, 1.9770199298095743e-21, -2.3406566479399704e-22, 1.713480149663986e-23}
	LANCZOS       = [...]float64{0.9999999999999971, 57.15623566586292, -59.59796035547549, 14.136097974741746, -0.4919138160976202, 3.399464998481189e-5, 4.652362892704858e-5, -9.837447530487956e-5, 1.580887032249125e-4, -2.1026444172410488e-4, 2.1743961811521265e-4, -1.643181065367639e-4, 8.441822398385275e-5, -2.6190838401581408e-5, 3.6899182659531625e-6}
	HALF_LOG_2_PI = 0.5 * math.Log(6.283185307179586)
)

type BinomialDistribution struct {
	trials      int64
	probability float64
	beta        *Beta
}

func NewBinomialDistribution(n int64, p float64) *BinomialDistribution {
	continuedFraction := &ContinuedFraction{}
	gamma := &Gamma{}
	beta := &Beta{
		defaultEpsilon: 1.0e-14,
		cf:             continuedFraction,
		gamma:          gamma,
	}
	return &BinomialDistribution{
		trials:      n,
		probability: p,
		beta:        beta,
	}
}

func (bd *BinomialDistribution) CumulativeProbability(x int64) (float64, error) {
	var ret float64
	if x < 0 {
		ret = 0.0
	} else if x >= bd.trials {
		ret = 1.0
	} else {
		value, err := bd.beta.SimpleRegularizedBeta(bd.probability, float64(x)+1.0, float64(bd.trials-x))
		if nil != err {
			return 0, err
		}
		ret = 1.0 - value
	}
	return ret, nil
}

func (bd *BinomialDistribution) InverseCumulativeProbability(p float64) (int64, error) {
	if p >= 0.0 && p <= 1.0 {
		lower := bd.getSupportLowerBound()
		if p == 0.0 {
			return lower, nil
		} else {
			if lower == -9223372036854775808 {
				if value, err := bd.checkedCumulativeProbability(lower); nil != err {
					return 0, err
				} else if value >= p {
					return lower, nil
				}
			} else {
				lower--
			}

			upper := bd.getSupportUpperBound()
			if p == 1.0 {
				return upper, nil
			} else {
				mu := bd.getNumericalMean()
				sigma := math.Sqrt(bd.getNumericalVariance())
				chebyshevApplies := !math.IsInf(mu, 0) && !math.IsNaN(mu) && !math.IsInf(sigma, 0) && !math.IsNaN(sigma) && sigma != 0.0
				if chebyshevApplies {
					k := math.Sqrt((1.0 - p) / p)
					tmp := mu - k*sigma
					if tmp > float64(lower) {
						lower = int64(math.Ceil(tmp)) - 1
					}

					k = 1.0 / k
					tmp = mu + k*sigma
					if tmp < float64(upper) {
						upper = int64(math.Ceil(tmp)) - 1
					}
				}
				return bd.solveInverseCumulativeProbability(p, lower, upper)
			}
		}
	} else {
		return 0, fmt.Errorf("%v out of [%v, %v] range", p, 0, 1)
	}
}

func (bd *BinomialDistribution) solveInverseCumulativeProbability(p float64, lower int64, upper int64) (int64, error) {
	for lower+1 < upper {
		xm := (lower + upper) / 2
		if xm < lower || xm > upper {
			xm = lower + (upper-lower)/2
		}

		pm, err := bd.checkedCumulativeProbability(xm)
		if nil != err {
			return 0, err
		}
		if pm >= p {
			upper = xm
		} else {
			lower = xm
		}
	}
	return upper, nil
}

func (bd *BinomialDistribution) getSupportLowerBound() int64 {
	if bd.probability < 1.0 {
		return 0
	}
	return bd.trials
}

func (bd *BinomialDistribution) getSupportUpperBound() int64 {
	if bd.probability > 0.0 {
		return bd.trials
	}
	return 0
}

func (bd *BinomialDistribution) checkedCumulativeProbability(argument int64) (float64, error) {
	result := 0.0
	if value, err := bd.CumulativeProbability(argument); nil != err {
		return 0, err
	} else {
		result = value
	}
	if math.IsNaN(result) {
		return 0, fmt.Errorf("Discrete cumulative probability function returned NaN for argument %v", argument)
	}
	return result, nil
}

func (bd *BinomialDistribution) getNumericalMean() float64 {
	//return float64(float64(bd.trials) * bd.probability)
	return float64(bd.trials) * bd.probability
}

func (bd *BinomialDistribution) getNumericalVariance() float64 {
	p := bd.probability
	//return float64(float64(bd.trials) * p * (1.0 - p))
	return float64(bd.trials) * p * (1.0 - p)
}

type Beta struct {
	defaultEpsilon float64
	cf             *ContinuedFraction
	gamma          *Gamma
}

func (beta *Beta) SimpleRegularizedBeta(x float64, a float64, b float64) (float64, error) {
	return beta.RegularizedBeta(x, a, b, 1.0e-14, 9223372036854775807)
}

func (beta *Beta) RegularizedBeta(x float64, a float64, b float64, epsilon float64, maxIterations int64) (float64, error) {
	ret := 0.0
	if !math.IsNaN(x) && !math.IsNaN(a) && !math.IsNaN(b) && x >= 0.0 && x <= 1.0 && a > 0.0 && b > 0.0 {
		if x > (a+1.0)/(2.0+b+a) && 1.0-x <= (b+1.0)/(2.0+b+a) {
			value, err := beta.RegularizedBeta(1.0-x, b, a, epsilon, maxIterations)
			if nil != err {
				return 0, err
			}
			ret = 1.0 - value
		} else {
			lbr, err := beta.logBeta(a, b)
			if nil != err {
				return 0, err
			}
			if value, err := beta.cf.evaluate(a, b, x, epsilon, maxIterations); nil != err {
				return 0, err
			} else {
				ret = math.Exp(a*math.Log(x)+b*math.Log1p(-x)-math.Log(a)-lbr) * 1.0 / value
			}
		}
	}
	return ret, nil
}

func (beta *Beta) logGammaMinusLogGammaSum(a float64, b float64) (float64, error) {
	if a < 0.0 {
		return 0, numberIsTooSmallException(a, 0, true)
	} else if b < 10.0 {
		return 0, numberIsTooSmallException(b, 10, true)
	} else {
		var d float64
		var w float64
		if a <= b {
			d = b + (a - 0.5)
			if value, err := beta.deltaMinusDeltaSum(a, b); nil != err {
				return 0, err
			} else {
				w = value
			}
		} else {
			d = a + (b - 0.5)
			if value, err := beta.deltaMinusDeltaSum(b, a); nil != err {
				return 0, err
			} else {
				w = value
			}
		}

		u := d * math.Log1p(a/b)
		v := a * (math.Log(b) - 1.0)
		if u <= v {
			return w - u - v, nil
		} else {
			return w - v - u, nil
		}
	}
}

func (beta *Beta) sumDeltaMinusDeltaSum(p float64, q float64) (float64, error) {
	if p < 10.0 {
		return 0, numberIsTooSmallException(p, 10, true)
	} else if q < 10.0 {
		return 0, numberIsTooSmallException(q, 10, true)
	} else {
		a := math.Min(p, q)
		b := math.Max(p, q)
		sqrtT := 10.0 / a
		t := sqrtT * sqrtT
		z := DELTA[len(DELTA)-1]

		for i := len(DELTA) - 2; i >= 0; i-- {
			z = t*z + DELTA[i]
		}

		if value, err := beta.deltaMinusDeltaSum(a, b); nil != err {
			return 0, err
		} else {
			return z/a + value, nil
		}
	}
}

func (beta *Beta) deltaMinusDeltaSum(a float64, b float64) (float64, error) {
	if a >= 0.0 && a <= b {
		if b < 10.0 {
			return 0, numberIsTooSmallException(b, 10, true)
		} else {
			h := a / b
			p := h / (1.0 + h)
			q := 1.0 / (1.0 + h)
			q2 := q * q
			s := make([]float64, len(DELTA))
			s[0] = 1.0

			for i := 1; i < len(s); i++ {
				s[i] = float64(1.0) + float64(q) + float64(q2)*s[i-1]
			}

			sqrtT := 10.0 / b
			t := sqrtT * sqrtT
			w := DELTA[len(DELTA)-1] * s[len(s)-1]

			for i := len(DELTA) - 2; i >= 0; i-- {
				w = float64(t)*w + DELTA[i]*s[i]
			}

			//return float64(w) * float64(p) / float64(b), nil
			return float64(w) * p / b, nil
		}
	} else {
		return 0, fmt.Errorf("%v out of [%v, %v] range", a, 0, b)
	}
}

func (beta *Beta) logGammaSum(a float64, b float64) (float64, error) {
	if a >= 1.0 && a <= 2.0 {
		if b >= 1.0 && b <= 2.0 {
			x := a - 1.0 + (b - 1.0)
			if x <= 0.5 {
				return beta.gamma.logGamma1p(1.0 + x)
			} else {
				if x <= 1.5 {
					if value, err := beta.gamma.logGamma1p(x); nil != err {
						return 0, err
					} else {
						return value + math.Log1p(x), nil
					}
				} else {
					if value, err := beta.gamma.logGamma1p(x - 1.0); nil != err {
						return 0, err
					} else {
						return value + math.Log(x*(1.0+x)), nil
					}
				}
			}
		} else {
			return 0, errors.New(fmt.Sprintf("%v out of [%v, %v] range", b, 1, 2))
		}
	} else {
		return 0, errors.New(fmt.Sprintf("%v out of [%v, %v] range", a, 1, 2))
	}
}

func (beta *Beta) logBeta(p float64, q float64) (float64, error) {
	if !math.IsNaN(p) && !math.IsNaN(q) && p > 0.0 && q > 0.0 {
		a := math.Min(p, q)
		b := math.Max(p, q)
		var prod1 float64
		var ared float64
		var prod2 float64
		var bred float64
		if a >= 10.0 {
			if value, err := beta.sumDeltaMinusDeltaSum(a, b); nil != err {
				return 0, err
			} else {
				prod1 = value
			}
			ared = a / b
			prod2 = ared / (1.0 + ared)
			bred = -(a - 0.5) * math.Log(prod2)
			v := b * math.Log1p(ared)
			if bred <= v {
				return -0.5*math.Log(b) + 0.9189385332046727 + prod1 - bred - v, nil
			} else {
				return -0.5*math.Log(b) + 0.9189385332046727 + prod1 - v - bred, nil
			}
		} else if a > 2.0 {
			if b > 1000.0 {
				n := int64(math.Floor(a - 1.0))
				prod := 1.0
				ared := a

				var i int64
				for i = 0; i < n; i++ {
					ared--
					prod *= ared / (1.0 + ared/b)
				}

				gammaSumValue, err := beta.logGammaMinusLogGammaSum(ared, b)
				if nil != err {
					return 0, err
				}
				if value, err := beta.gamma.logGamma(ared); nil != err {
					return 0, err
				} else {
					return math.Log(prod) - float64(n)*math.Log(b) + value + gammaSumValue, nil
				}
			} else {
				prod1 = 1.0

				for ared = a; ared > 2.0; prod1 *= prod2 / (1.0 + prod2) {
					ared--
					prod2 = ared / b
				}

				if b >= 10.0 {
					gammaSumValue, err := beta.logGammaMinusLogGammaSum(ared, b)
					if nil != err {
						return 0, err
					}
					if value, err := beta.gamma.logGamma(ared); nil != err {
						return 0, err
					} else {
						return math.Log(prod1) + value + gammaSumValue, nil
					}
				} else {
					prod2 = 1.0

					for bred = b; bred > 2.0; prod2 *= bred / (ared + bred) {
						bred--
					}

					if value, err := beta.logGammaSum(ared, bred); nil != err {
						return 0, err
					} else {
						value1, err := beta.gamma.logGamma(ared)
						if nil != err {
							return 0, err
						}
						value2, err := beta.gamma.logGamma(bred)
						if nil != err {
							return 0, err
						}
						return math.Log(prod1) + math.Log(prod2) + value1 + (value2 - value), nil
					}
				}
			}
		} else if a < 1.0 {
			if b >= 10.0 {
				gammaSumValue, err := beta.logGammaMinusLogGammaSum(a, b)
				if nil != err {
					return 0, err
				}
				if value, err := beta.gamma.logGamma(a); nil != err {
					return 0, err
				} else {
					return value + gammaSumValue, nil
				}
			} else {
				v1, err := beta.gamma.gamma(a)
				if nil != err {
					return 0, err
				}
				v2, err := beta.gamma.gamma(b)
				if nil != err {
					return 0, err
				}
				v3, err := beta.gamma.gamma(a + b)
				if nil != err {
					return 0, err
				}
				return math.Log(v1 * v2 / v3), nil
			}
		} else if b <= 2.0 {
			if value, err := beta.logGammaSum(a, b); nil != err {
				return 0, err
			} else {
				lgv1, err := beta.gamma.logGamma(a)
				if nil != err {
					return 0, err
				}
				lgv2, err := beta.gamma.logGamma(b)
				if nil != err {
					return 0, err
				}
				return lgv1 + lgv2 - value, nil
			}
		} else if b >= 10.0 {
			gammaSumValue, err := beta.logGammaMinusLogGammaSum(a, b)
			if nil != err {
				return 0, err
			}
			if value, err := beta.gamma.logGamma(a); nil != err {
				return 0, err
			} else {
				return value + gammaSumValue, nil
			}
		} else {
			prod1 = 1.0

			for ared = b; ared > 2.0; prod1 *= ared / (a + ared) {
				ared--
			}

			if value, err := beta.logGammaSum(a, ared); nil != err {
				return 0, err
			} else {
				lgv1, err := beta.gamma.logGamma(a)
				if nil != err {
					return 0, err
				}
				lgv2, err := beta.gamma.logGamma(ared)
				if nil != err {
					return 0, err
				}
				return math.Log(prod1) + lgv1 + (lgv2 - value), nil
			}
		}
	} else {
		return 0.0, nil
	}
}

type ContinuedFraction struct {
}

func (cf *ContinuedFraction) getA(n int64, x float64) float64 {
	return 1.0
}

func (cf *ContinuedFraction) getB(a float64, b float64, n int64, x float64) (ret float64) {
	var m float64
	if n%2 == 0 {
		m = float64(n) / 2.0
		ret = m * (b - m) * x / ((a + 2.0*m - 1.0) * (a + 2.0*m))
	} else {
		m = (float64(n) - 1.0) / 2.0
		ret = -((a + m) * (a + b + m) * x) / ((a + 2.0*m) * (a + 2.0*m + 1.0))
	}
	return ret
}

func (cf *ContinuedFraction) evaluate(av float64, bv float64, x float64, epsilon float64, maxIterations int64) (float64, error) {
	hPrev := cf.getA(0, x)
	if precisionEq(hPrev, 0.0, 1.0e-50) {
		hPrev = 1.0e-50
	}

	var n int64 = 1
	dPrev := 0.0
	cPrev := hPrev
	hN := hPrev

	for {
		if n < maxIterations {
			a := cf.getA(n, x)
			b := cf.getB(av, bv, n, x)
			dN := a + b*dPrev
			if precisionEq(dN, 0.0, 1.0e-50) {
				dN = 1.0e-50
			}

			cN := a + b/cPrev
			if precisionEq(cN, 0.0, 1.0e-50) {
				cN = 1.0e-50
			}

			dN = 1.0 / dN
			deltaN := cN * dN
			hN = hPrev * deltaN
			if math.IsInf(hN, 0) {
				return 0, errors.New(fmt.Sprintf("Continued fraction convergents diverged to +/- infinity for value %v", x))
			}

			if math.IsNaN(hN) {
				return 0, errors.New(fmt.Sprintf("Continued fraction diverged to NaN for value %v", x))
			}

			if math.Abs(deltaN-1.0) >= epsilon {
				dPrev = dN
				cPrev = cN
				hPrev = hN
				n++
				continue
			}
		}

		if n >= maxIterations {
			return 0, errors.New(fmt.Sprintf("Continued fraction convergents failed to converge (in less than %v iterations) for value %v", maxIterations, x))
		}

		return hN, nil
	}
}

func precisionEq(x float64, y float64, eps float64) bool {
	return precisionEqs(x, y, 1) || math.Abs(y-x) <= eps
}

func precisionEqs(x float64, y float64, maxUlps float64) bool {
	xInt := int64(math.Float64bits(x))
	yInt := int64(math.Float64bits(y))
	var isEqual bool
	if ((xInt ^ yInt) & -9223372036854775808) == 0 {
		isEqual = int64(math.Abs(float64(xInt)-float64(yInt))) <= int64(maxUlps)
	} else {
		var deltaPlus int64
		var deltaMinus int64
		if xInt < yInt {
			deltaPlus = yInt - int64(math.Float64bits(0.0))
			deltaMinus = xInt - int64(math.Float64bits(-0.0))
		} else {
			deltaPlus = xInt - int64(math.Float64bits(0.0))
			deltaMinus = yInt - int64(math.Float64bits(-0.0))
		}

		if deltaPlus > int64(maxUlps) {
			isEqual = false
		} else {
			isEqual = deltaMinus <= int64(maxUlps)-deltaPlus
		}
	}

	return isEqual && !math.IsNaN(x) && !math.IsNaN(y)
}

type Gamma struct {
}

func (g *Gamma) logGamma(x float64) (float64, error) {
	ret := 0.0
	if !math.IsNaN(x) && x > 0.0 {
		if x < 0.5 {
			if value, err := g.logGamma1p(x); nil != err {
				return 0, err
			} else {
				return value - math.Log(x), nil
			}
		}

		if x <= 2.5 {
			return g.logGamma1p(x - 0.5 - 0.5)
		}

		if x <= 8.0 {
			n := int64(math.Floor(x - 1.5))
			prod := 1.0

			var i int64
			for i = 1; i <= n; i++ {
				prod *= x - float64(i)
			}

			if value, err := g.logGamma1p(x - float64(n+1)); nil != err {
				return 0, err
			} else {
				return value + math.Log(prod), nil
			}
		}

		sum := g.lanczos(x)
		tmp := x + 4.7421875 + 0.5
		ret = (x+0.5)*math.Log(tmp) - tmp + HALF_LOG_2_PI + math.Log(sum/x)
	}
	return ret, nil
}

func (g *Gamma) logGamma1p(x float64) (float64, error) {
	if x < -0.5 {
		return 0, numberIsTooSmallException(x, -0.5, true)
	} else if x > 1.5 {
		return 0, numberIsTooLargeException(x, 1.5, true)
	} else {
		if value, err := g.invGamma1pm1(x); nil != err {
			return 0, err
		} else {
			return -math.Log1p(value), nil
		}
	}
}

func (g *Gamma) invGamma1pm1(x float64) (float64, error) {
	if x < -0.5 {
		return 0, numberIsTooSmallException(x, -0.5, true)
	} else if x > 1.5 {
		return 0, numberIsTooLargeException(x, 1.5, true)
	} else {
		t := x
		if x > 0.5 {
			t = x - float64(0.5) - float64(0.5)
		}
		var ret float64
		var a float64
		var b float64
		var c float64
		if t < 0.0 {
			a = 6.116095104481416e-9 + t*6.247308301164655e-9
			b = 1.9575583661463974e-10
			b = -6.077618957228252e-8 + t*b
			b = 9.926418406727737e-7 + t*b
			b = -6.4304548177935305e-6 + t*b
			b = -8.514194324403149e-6 + t*b
			b = 4.939449793824468e-4 + t*b
			b = 0.026620534842894922 + t*b
			b = 0.203610414066807 + t*b
			b = 1.0 + t*b
			c = -2.056338416977607e-7 + t*(a/b)
			c = 1.133027231981696e-6 + t*c
			c = -1.2504934821426706e-6 + t*c
			c = -2.013485478078824e-5 + t*c
			c = 1.280502823881162e-4 + t*c
			c = -2.1524167411495098e-4 + t*c
			c = -0.0011651675918590652 + t*c
			c = 0.0072189432466631 + t*c
			c = -0.009621971527876973 + t*c
			c = -0.04219773455554433 + t*c
			c = 0.16653861138229148 + t*c
			c = -0.04200263503409524 + t*c
			c = -0.6558780715202539 + t*c
			c = -0.42278433509846713 + t*c
			if x > 0.5 {
				ret = t * c / x
			} else {
				ret = x * (c + 0.5 + 0.5)
			}
		} else {
			a = 4.343529937408594e-15
			a = -1.2494415722763663e-13 + t*a
			a = 1.5728330277104463e-12 + t*a
			a = 4.686843322948848e-11 + t*a
			a = 6.820161668496171e-10 + t*a
			a = 6.8716741130671986e-9 + t*a
			a = 6.116095104481416e-9 + t*a
			b = 2.6923694661863613e-4
			b = 0.004956830093825887 + t*b
			b = 0.054642130860422966 + t*b
			b = 0.3056961078365221 + t*b
			b = 1.0 + t*b
			c = -2.056338416977607e-7 + a/b*t
			c = 1.133027231981696e-6 + t*c
			c = -1.2504934821426706e-6 + t*c
			c = -2.013485478078824e-5 + t*c
			c = 1.280502823881162e-4 + t*c
			c = -2.1524167411495098e-4 + t*c
			c = -0.0011651675918590652 + t*c
			c = 0.0072189432466631 + t*c
			c = -0.009621971527876973 + t*c
			c = -0.04219773455554433 + t*c
			c = 0.16653861138229148 + t*c
			c = -0.04200263503409524 + t*c
			c = -0.6558780715202539 + t*c
			c = 0.5772156649015329 + t*c
			if x > 0.5 {
				ret = t / x * (c - 0.5 - 0.5)
			} else {
				ret = x * c
			}
		}

		return ret, nil
	}
}

func (g *Gamma) lanczos(x float64) float64 {
	sum := 0.0

	for i := len(LANCZOS) - 1; i > 0; i-- {
		sum += LANCZOS[i] / (x + float64(i))
	}

	return sum + LANCZOS[0]
}

func (g *Gamma) gamma(x float64) (float64, error) {
	if x == math.Floor(x+0.5) && x <= 0.0 {
		return 0.0, nil
	} else {
		var ret float64
		absX := math.Abs(x)
		var prod float64
		var t float64
		if absX <= 20.0 {
			if x >= 1.0 {
				prod = 1.0

				for t := x; t > 2.5; prod *= t {
					t--
				}

				if value, err := g.invGamma1pm1(t - 1.0); nil != err {
					return 0, err
				} else {
					ret = prod / (1.0 + value)
				}
			} else {
				prod = x

				for t := x; t < -0.5; prod *= t {
					t++
				}

				if value, err := g.invGamma1pm1(t); nil != err {
					return 0, err
				} else {
					ret = 1.0 / (prod * (1.0 + value))
				}
			}
		} else {
			prod = absX + 4.7421875 + 0.5
			t = 2.5066282746310007 / absX * math.Pow(prod, absX+0.5) * math.Exp(-prod) * g.lanczos(absX)
			if x > 0.0 {
				ret = t
			} else {
				ret = -3.141592653589793 / (x * math.Sin(3.141592653589793*x) * t)
			}
		}

		return ret, nil
	}
}

func numberIsTooSmallException(wrong float64, min float64, boundIsAllowed bool) error {
	if boundIsAllowed {
		return errors.New(fmt.Sprintf("%v is smaller than the minimum (%v)", wrong, min))
	}
	return errors.New(fmt.Sprintf("%v is smaller than, or equal to, the minimum (%v)", wrong, min))
}
func numberIsTooLargeException(wrong float64, min float64, boundIsAllowed bool) error {
	if boundIsAllowed {
		return errors.New(fmt.Sprintf("%v is larger than the maximum (%v)", wrong, min))
	}
	return errors.New(fmt.Sprintf("%v is larger than, or equal to, the maximum (%v)", wrong, min))
}
