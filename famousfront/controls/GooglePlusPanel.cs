using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using famousfront.utils;
namespace famousfront.controls
{
  public class GooglePlusPanel : Panel
  {
    public class PanelOptions
    {
      public int ColumnCount = 2;
      public int DefaultColumnWidth = 320;
      public double LineHeight = 30d;
    }
    public PanelOptions Options { get; set; }

    public GooglePlusPanel()
    {
      Options = new PanelOptions();
    }

    protected override Size MeasureOverride(Size avail)
    {
      var sz = avail;
      if (double.IsInfinity(sz.Width))
        sz.Width = Options.ColumnCount * Options.DefaultColumnWidth;
      var colwidth = sz.Width / Options.ColumnCount;
      var colheights = new double[Options.ColumnCount];
      var enable_span = new bool[Options.ColumnCount];
      var lastitems = new UIElement[Options.ColumnCount];

      foreach (UIElement child in InternalChildren)
      {
        var col = child.GetColumnIndex();
        int span_count = child.GetColumnSpan();
        double offset;
        if (col < 0)
          col = select_column(colheights, enable_span, out span_count, out offset);
        else
        {
          offset = colheights[col];
        }
        var itemsz = new Size(colwidth * (span_count + 1), double.PositiveInfinity);
        child.Measure(itemsz);
        var rect = new Rect(col * colwidth, offset, colwidth * (span_count + 1), child.DesiredSize.Height);
        child.SetPosition(rect);

        update_colheight(colheights, col, span_count, child.DesiredSize.Height, enable_span, offset);
        update_lastitem_layout(lastitems, col, span_count, offset, child);
      }
      sz.Height = columns_height(colheights);
      return sz;
    }
    void update_lastitem_layout(IList<UIElement> lastitems, int col, int span_count, double offset, UIElement child)
    {
      if (span_count > 0)
      {
        for (var i = 0; i <= span_count; i++)
        {
          var item = lastitems[col + i];
          if (item == null)
          {
            continue;
          }
          var rc = item.GetPosition();
          item.SetPosition(new Rect(rc.Left, rc.Top, rc.Width, offset - rc.Top));
        }
      }

      lastitems[col] = child;
      for (var i = 1; i <= span_count; i++)
      {
        lastitems[col + i] = null;
      }

    }

    static double columns_height(IEnumerable<double> columns)
    {
      return columns.Max();
    }
    protected override Size ArrangeOverride(Size finalsz)
    {
      foreach (UIElement child in InternalChildren)
      {
        var pos = child.GetPosition();
        child.Arrange(pos);
      }
      return finalsz;
    }

    private void update_colheight(double[] columns, int col_start, int span_count, double height, bool[] enable_span, double offset)
    {
      for (var i = 0; i <= span_count; ++i)
      {
        columns[col_start + i] = offset + height;
      }
      enable_span[col_start] = span_count == 0;
    }
    private int select_column(IList<double> columns, bool[] enable_span, out int span_count, out double offset)
    {
      span_count = 0;
      var min = double.PositiveInfinity;
      var idx = 0;
      var cnt = columns.Count();
      for (var i = 0; i < cnt; ++i)
      {
        var h = Math.Ceiling(columns[i] / Options.LineHeight);
        if (h.lt(min))
        {
          idx = i;
          min = h;
        }
      }
      offset = columns[idx];

      if (!enable_span[idx])
        return idx;
      for (var i = idx + 1; i < cnt; ++i)
      {
        var h = Math.Ceiling(columns[i] / Options.LineHeight);
        if (h.eq(min))
        {
          offset = Math.Max(offset, columns[i]);
          ++span_count;
        }
        else
        {
          break;
        }
      }
      return idx;
    }
  }

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
