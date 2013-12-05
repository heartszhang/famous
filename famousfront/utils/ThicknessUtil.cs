using System.Diagnostics;
using System.Diagnostics.Contracts;
using System.Windows;

namespace famousfront.utils
{
  [DebuggerNonUserCode]
  [System.Diagnostics.Contracts.Pure]
  internal static class ThicknessUtil
  {
    internal static bool IsNonNegative(Thickness value)
    {
      return DoubleExtension.IsNonNegative(value.Left) &&
             DoubleExtension.IsNonNegative(value.Top) &&
             DoubleExtension.IsNonNegative(value.Right) &&
             DoubleExtension.IsNonNegative(value.Bottom);
    }

    [Conditional("CONTRACTS_FULL")]
    [ContractAbbreviator]
    internal static void EnsureNonNegative()
    {
      Contract.Ensures(IsNonNegative(Contract.Result<Thickness>()));
    }

    internal static object CoerceNonNegative(DependencyObject obj, object basevalue)
    {
      Debug.Assert(obj != null);
      var value = (Thickness)basevalue;
      if (!DoubleExtension.IsNonNegative(value.Left))
      {
        value.Left = 0d;
      }
      if (!DoubleExtension.IsNonNegative(value.Top))
      {
        value.Top = 0d;
      }
      if (!DoubleExtension.IsNonNegative(value.Right))
      {
        value.Right = 0d;
      }
      if (!DoubleExtension.IsNonNegative(value.Bottom))
      {
        value.Bottom = 0d;
      }
      return value;
    }
  }
}