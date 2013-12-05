﻿using System;
using System.Globalization;
using System.Windows;
using System.Windows.Data;

namespace famousfront.converters
{
  [ValueConversion(typeof(int), typeof(Visibility))]
  public sealed class IntToVisibilityConverter : IValueConverter
  {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture)
    {
      if (value == null)
        return Visibility.Collapsed;

      var count = (int)value;
      if (parameter != null)
      {
        return count == 0 ? Visibility.Visible : Visibility.Collapsed;
      }
      return count > 0 ? Visibility.Visible : Visibility.Collapsed;
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture)
    {
      throw new NotImplementedException();
    }
  }
}