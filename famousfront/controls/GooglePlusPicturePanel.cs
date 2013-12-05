using System.Collections.Generic;
using System.Windows;
using System.Windows.Controls;
using famousfront.utils;
namespace famousfront.controls
{
  public class GooglePlusPicturePanel : Panel
  {
    public GooglePlusPicturePanel()
    {
      WidthThreshold = 320;
      HeightThreshold = 240;
    }
    public static readonly DependencyProperty IsExpandedProperty = DependencyProperty.Register(
      "IsExpanded",      typeof(bool),      typeof(GooglePlusPicturePanel),
      new FrameworkPropertyMetadata(false,
          FrameworkPropertyMetadataOptions.AffectsRender | 
          FrameworkPropertyMetadataOptions.AffectsMeasure )    );
    public bool IsExpanded
    {
      get { return (bool)GetValue(IsExpandedProperty); }
      set { SetValue(IsExpandedProperty, value); }
    }
    public double WidthThreshold { get; set; }
    public double HeightThreshold { get; set; }
    protected override Size MeasureOverride(Size avail)
    {
      var mheight = 0.0;
      var unit_width = 0.0;
      var row = new List<UIElement>();
      var first_line = 0d;
      foreach (UIElement child in InternalChildren)
      {
        row.Add(child);
        var scale = child.GetScale();
        if (scale.lt(0.0))
        {
          child.Measure(new Size(avail.Width, HeightThreshold));
          scale = child.DesiredSize.Height.zero() ? -1.0 : (child.DesiredSize.Width / child.DesiredSize.Height);
          child.SetScale(scale);
        }
        unit_width += scale;
        if (HeightThreshold * unit_width > avail.Width)
        {
          var cheight = avail.Width / unit_width;
          do_row_measure(0, mheight, cheight, row);
          mheight += cheight;
          unit_width = 0.0;
          row.Clear();
          if (first_line.zero())
          {
            first_line = cheight;
          }
        }
      }
      if (row.Count > 0)
      {
        do_row_measure(0, mheight, HeightThreshold, row);
        mheight += HeightThreshold;
      }
      if (first_line.zero())
        first_line = HeightThreshold;
      if (!IsExpanded)
      {
        mheight = first_line;
      }
      return new Size(avail.Width, mheight);
    }

    void do_row_measure(double x, double y, double row_height, IEnumerable<UIElement> row)
    {
      foreach (var ue in row)
      {
        var scale = ue.GetScale();
        var w = row_height * scale;
        var pos = new Rect(x, y, w, row_height);
        x += w;
        ue.SetPosition(pos);
      }
    }
    protected override Size ArrangeOverride(Size finalsz)
    {
      foreach (UIElement child in InternalChildren)
      {
        var rect = child.GetPosition();
        child.Arrange(rect);
      }

      return finalsz;
    }
  }
}
