using System;
using System.Globalization;
using System.Windows;
using System.Windows.Data;

namespace famousfront.converters
{
  public sealed class BoolToMaxHeightConverter : IValueConverter
  {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture)
    {
      if (value == null)
        return parameter;

      var boolean = (bool)value;
      return boolean ? DependencyProperty.UnsetValue : parameter;
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture)
    {
      throw new NotImplementedException();
    }
  }
}