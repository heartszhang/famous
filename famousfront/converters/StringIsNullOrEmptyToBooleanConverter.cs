﻿using System;
using System.Globalization;
using System.Windows.Data;

namespace famousfront.converters
{
  [ValueConversion(typeof(string), typeof(bool))]
  public sealed class StringIsNullOrEmptyToBooleanConverter : IValueConverter
  {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture)
    {
      var str = value as string;

      return string.IsNullOrEmpty(str);
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture)
    {
      throw new NotImplementedException();
    }
  }
}