using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.utils
{
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
    }
}
