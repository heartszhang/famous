using System;
using famousfront.Properties;

namespace famousfront.utils
{
  internal class FriendlyDateTime 
  {
    DateTime _;
    internal FriendlyDateTime(DateTime t)
    {
      _ = t;
    }
    public static implicit operator DateTime(FriendlyDateTime t)
    {
      return t._;
    }

    public override string ToString()
    {
      var p = _;
      var now = DateTime.Now;
      var v = p.ToString("D", new System.Globalization.CultureInfo("zh-cn"));
      if (p.Year != now.Year)
        return v;
      var diff = (now - p).Days;
      var dw = (int)p.DayOfWeek - 1;
      var ndw = (int)now.DayOfWeek -1;

      if (dw < 0)
        dw = 6;      
      if (ndw < 0)
        ndw = 6;

      var firstdthisweek = now.AddDays(-ndw);
      var prevweek = firstdthisweek.AddDays(-7d);
      var ns = new[] { Resources.Today, Resources.Yesterday, Resources.DayBeforeYeserday, Resources.ThreeDaysAgo };
      var cws = new[] { Resources.Monday, Resources.Tuesday, Resources.Wednesday, Resources.Thusday, Resources.Friday, Resources.Saturday, Resources.Sunday };
      if (diff >= 0 && diff < ns.Length)
      {
        v = ns[diff];
      }else if (p >= firstdthisweek)
      {
        v = cws[dw];
      }
      else if (p >= prevweek)
      {
        v = "上" + cws[dw];
      }
      return v;
    }
  }
}