using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Media;

namespace famousfront.utils
{
    internal static class ColorExtensions
    {
        public static string ToHexValue(this Color c)
        {
            return String.Format("#{0}{1}{2}{3}", c.A.ToString("X2"), c.R.ToString("X2"), c.G.ToString("X2"), c.B.ToString("X2"));
        }

        public static Color ColorFromHexValue(string hexValue)
        {
            byte alpha;
            byte pos = 0;

            var hex = hexValue.Replace("#", "");

            if (hex.Length == 8)
            {
                alpha = Convert.ToByte(hex.Substring(pos, 2), 16);
                pos = 2;
            }
            else
            {
                alpha = Convert.ToByte("ff", 16);
            }

            var red = Convert.ToByte(hex.Substring(pos, 2), 16);

            pos += 2;
            var green = Convert.ToByte(hex.Substring(pos, 2), 16);

            pos += 2;
            var blue = Convert.ToByte(hex.Substring(pos, 2), 16);

            return Color.FromArgb(alpha, red, green, blue);
        }
    }
}
