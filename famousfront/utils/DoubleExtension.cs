using System;
using System.Diagnostics;
using System.Diagnostics.Contracts;
using System.Windows;

namespace famousfront.utils
{
  [DebuggerNonUserCode]
  [System.Diagnostics.Contracts.Pure]
  internal static class DoubleExtension
  {
    internal const double Epsilon = 0.000001;
    internal static bool gt(this double self, double rhs)
    {
      return (self - rhs) > Epsilon;
    }

    internal static bool lt(this double self, double rhs)
    {
      return (self - rhs) < -Epsilon;
    }

    internal static bool eq(this double self, double rhs)
    {
      return Math.Abs(self - rhs) < Epsilon;
    }

    internal static bool let(this double self, double rhs)
    {
      return (self - rhs) < Epsilon;
    }

    internal static bool get(this double self, double rhs)
    {
      return (self - rhs) > -Epsilon;
    }

    internal static bool zero(this double self)
    {
      return Math.Abs(self) < Epsilon;
    }
    internal static object CoerceNonNegative(DependencyObject obj, object basevalue)
    {
      Debug.Assert(obj != null);
      var value = (double)basevalue;
      return IsNonNegative(value) ? value : 0d;
    }
    [Conditional("CONTRACTS_FULL")]
    [ContractAbbreviator]
    internal static void EnsureNonNegative()
    {
      Contract.Ensures(condition: IsNonNegative(Contract.Result<double>()));
    }
    internal static bool IsNonNegative(double value)
    {
      return !double.IsNaN(value) && !double.IsInfinity(value) && value > 0d;
    }
  }
}
