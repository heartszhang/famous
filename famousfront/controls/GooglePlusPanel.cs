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
        public GooglePlusPanel()
        {
            WidthThreshold = 320;
            HeightThreshold = 240;
        }
        public double WidthThreshold { get; set; }
        public double HeightThreshold { get; set; }
        protected override Size MeasureOverride(Size avail)
        {
            var mheight = 0.0;
            var unit_width = 0.0;
            var row = new List<UIElement>();
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
                }
            }
            if (row.Count > 0)
            {
                do_row_measure(0, mheight, HeightThreshold, row);
                mheight += HeightThreshold;
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
    internal static class GooglePlusPanelProperties
    {
        public static readonly DependencyProperty ScaleProperty = DependencyProperty.RegisterAttached("Scale",
            typeof(double), typeof(GooglePlusPanelProperties), new PropertyMetadata(-1.0));

        public static readonly DependencyProperty PositionProperty = DependencyProperty.RegisterAttached("Position",
            typeof(Rect), typeof(GooglePlusPanelProperties));

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
