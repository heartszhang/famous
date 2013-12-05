using System;
using System.Globalization;
using System.Windows;
using System.Windows.Data;

namespace famousfront.converters
{
  [ValueConversion(typeof(string), typeof(Visibility))]
  public sealed class StringIsNullOrEmptyToVisibilityConverter : IValueConverter
  {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture)
    {
      var str = value as string;
      var en = string.IsNullOrEmpty(str);
      if (parameter != null)
        en = !en;
      return en ? Visibility.Collapsed : Visibility.Visible;
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture)
    {
      throw new NotImplementedException();
    }
  }
}