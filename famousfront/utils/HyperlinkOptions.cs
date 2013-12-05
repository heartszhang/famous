using System.Diagnostics;
using System.Windows;
using System.Windows.Media;

namespace famousfront.utils
{
  public static class HyperlinkOptions
  {
    public static readonly DependencyProperty ForegroundBrushProperty =
        DependencyProperty.RegisterAttached("ForegroundBrush", typeof(SolidColorBrush), typeof(HyperlinkOptions),
                                            new FrameworkPropertyMetadata(null,
                                                                          FrameworkPropertyMetadataOptions.AffectsRender |
                                                                          FrameworkPropertyMetadataOptions.Inherits));


    public static SolidColorBrush GetForegroundBrush(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return obj.GetValue(ForegroundBrushProperty) as SolidColorBrush;
    }


    public static void SetForegroundBrush(DependencyObject obj, SolidColorBrush value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(ForegroundBrushProperty, value);
    }
  }
}
