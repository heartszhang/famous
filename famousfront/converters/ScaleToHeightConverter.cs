using System;
using System.Globalization;
using System.Windows.Data;
using famousfront.utils;

namespace famousfront.converters
{
  [ValueConversion(typeof(double), typeof(double))]
  public sealed class ScaleToHeightConverter : IValueConverter
  {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture)
    {
      var width = System.Convert.ToDouble(parameter);
      var scale = System.Convert.ToDouble(value);
      return scale.zero() || width.let(0.0) || scale.let(0.0) ? 0 : width / scale;
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture)
    {
      throw new NotImplementedException();
    }
  }
}