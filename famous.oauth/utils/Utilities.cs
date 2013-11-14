using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;

namespace famous.oauth.utils
{
    /// <summary>A utility class which contains helper methods and extension methods.</summary>
    internal static class Utilities
    {
        /// <summary>Returns the version of the core library.</summary>
        internal static string GetLibraryVersion()
        {
            return Regex.Match(typeof(Utilities).Assembly.FullName, "Version=([\\d\\.]+)").Groups[1].ToString();
        }       
    }
}
