using System;
using System.Globalization;
using System.Windows.Data;

namespace famousfront.converters
{
  [ValueConversion(typeof(bool), typeof(Elysium.Controls.ProgressState))]
  internal sealed class BoolToBusyIndicatorConverter : IValueConverter
  {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture)
    {
      if (value == null)
        return Elysium.Controls.ProgressState.Normal;

      var boolean = (bool)value;

      if (parameter != null)
        boolean = !boolean;

      return boolean ? Elysium.Controls.ProgressState.Indeterminate : Elysium.Controls.ProgressState.Normal;
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture)
    {
      throw new NotImplementedException();
    }
  }
}