using System.Windows;
using System.Windows.Controls;
using famousfront.utils;

namespace famousfront.controls
{
  public class ExpandableStackPanel : StackPanel
  {
    public static readonly DependencyProperty CanExpandProperty =
      DependencyProperty.Register("CanExpand", typeof(bool),typeof(ExpandableStackPanel));

    public static readonly DependencyProperty IsExpandedProperty =
      DependencyProperty.Register("IsExpanded", typeof(bool),
        typeof(ExpandableStackPanel), new FrameworkPropertyMetadata(false, 
          FrameworkPropertyMetadataOptions.AffectsMeasure | FrameworkPropertyMetadataOptions.AffectsArrange));

    public static readonly DependencyProperty MiniHeightProperty =
      DependencyProperty.Register("MiniHeight", typeof(double),
        typeof(ExpandableStackPanel), new FrameworkPropertyMetadata(412d, 
          FrameworkPropertyMetadataOptions.AffectsMeasure));

    public bool CanExpand
    {
      get { return (bool)GetValue(CanExpandProperty); }
      set { SetValue(CanExpandProperty, value); }
    }

    public bool IsExpanded
    {
      get { return (bool)GetValue(IsExpandedProperty); }
      set { SetValue(IsExpandedProperty, value); }
    }
    public double MiniHeight
    {
      get { return (double)GetValue(MiniHeightProperty); }
      set { SetValue(MiniHeightProperty, value); }
    }
    protected override Size MeasureOverride(Size constraint)
    {
      var sz = base.MeasureOverride(constraint);
      CanExpand = sz.Height.gt(MiniHeight);
      if (!CanExpand)
      {
        return sz;
      }
      if (!IsExpanded)
      {
        sz.Height = MiniHeight;
      }
      return sz;
    }
  }
}