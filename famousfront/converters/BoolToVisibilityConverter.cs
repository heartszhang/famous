using System;
using System.Globalization;
using System.Windows;
using System.Windows.Data;

namespace famousfront.converters
{
  [ValueConversion(typeof(bool), typeof(Visibility))]
  internal class BoolToVisibilityConverter : IValueConverter
  {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture)
    {
      if (value == null)
        return Visibility.Collapsed;

      var boolean = (bool)value;

      if (parameter != null)
        boolean = !boolean;

      return boolean ? Visibility.Visible : Visibility.Collapsed;
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture)
    {
      throw new NotImplementedException();
    }
  }
}
