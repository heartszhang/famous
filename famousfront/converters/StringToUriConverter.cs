using System;
using System.Globalization;
using System.Windows.Data;

namespace famousfront.converters
{
  [ValueConversion(typeof(string), typeof(Uri))]
  public sealed class StringToUriConverter : IValueConverter
  {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture)
    {
      var str = value as string;
      var en = string.IsNullOrEmpty(str);
      return string.IsNullOrEmpty(str) ? null : new Uri(str);
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture)
    {
      throw new NotImplementedException();
    }
  }
}