using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.utils
{
    class UrlLauncher
    {
        public void LaunchUri(Uri location)
        {
            if (location == null)
                throw new ArgumentNullException("location");

            // ensure it's either http/s or mailto
            var scheme = location.Scheme.ToLowerInvariant();
            if (scheme == "http" ||
                scheme == "https" ||
                scheme == "mailto")
            {
                var str = location.ToString();

                using (var p = Process.Start(str)) { }
            }
        }
    }
}
