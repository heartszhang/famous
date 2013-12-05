using System.Windows;

namespace famousfront.controls
{
  internal static class LayoutProperties
  {
    public static readonly DependencyProperty ScaleProperty = DependencyProperty.RegisterAttached("Scale",
      typeof(double), typeof(LayoutProperties), new PropertyMetadata(-1.0));

    public static readonly DependencyProperty PositionProperty = DependencyProperty.RegisterAttached("Position",
      typeof(Rect), typeof(LayoutProperties));

    public static readonly DependencyProperty ColumnSpanProperty = DependencyProperty.RegisterAttached("ColumnSpan",
      typeof(int), typeof(LayoutProperties));

    public static void SetColumnSpan(this UIElement ue, int val)
    {
      ue.SetValue(ColumnSpanProperty, val);
    }

    public static int GetColumnSpan(this UIElement ue)
    {
      return (int)ue.GetValue(ColumnSpanProperty);
    }

    public static readonly DependencyProperty ColumnIndexProperty = DependencyProperty.RegisterAttached("ColumnIndex",
      typeof(int), typeof(LayoutProperties), new PropertyMetadata(-1));

    public static void SetColumnIndex(this UIElement ue, int val)
    {
      ue.SetValue(ColumnIndexProperty, val);
    }

    public static int GetColumnIndex(this UIElement ue)
    {
      return (int)ue.GetValue(ColumnIndexProperty);
    }

    public static void SetScale(this UIElement ue, double val)
    {
      ue.SetValue(ScaleProperty, val);
    }

    public static double GetScale(this UIElement ue)
    {
      return (double)ue.GetValue(ScaleProperty);
    }

    public static void SetPosition(this UIElement ue, Rect val)
    {
      ue.SetValue(PositionProperty, val);
    }

    public static Rect GetPosition(this UIElement ue)
    {
      return (Rect)ue.GetValue(PositionProperty);
    }
  }
}